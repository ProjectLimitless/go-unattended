/**
* This file is part of Unattended.
* Copyright © 2018 Donovan Solms.
* Project Limitless
* https://www.projectlimitless.io
*
* Unattended and Project Limitless is free software: you can redistribute it and/or modify
* it under the terms of the Apache License Version 2.0.
*
* You should have received a copy of the Apache License Version 2.0 with
* Unattended. If not, see http://www.apache.org/licenses/LICENSE-2.0.
 */

package unattended

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/ProjectLimitless/go-unattended/omaha"
	"github.com/cavaliercoder/grab"
	"github.com/otiai10/copy"
	"github.com/sirupsen/logrus"
)

// Unattended implements the core functionality of the package. It takes
// ownership of running and updating a target application
type Unattended struct {
	mutex               sync.Mutex
	clientID            string
	target              Target
	updateCheckInterval time.Duration
	outputWriter        io.Writer
	// command holds the target application when executed
	command          *exec.Cmd
	commandCompleted bool
	log              *logrus.Entry
	waitGroup        sync.WaitGroup
}

// New creates a new instance of the unattended updater
func New(
	clientID string,
	target Target,
	updateCheckInterval time.Duration,
	log *logrus.Entry) (*Unattended, error) {

	if target.VersionsPath == "" || target.ApplicationName == "/" {
		return nil, fmt.Errorf(
			"Target version path '%s' is not valid",
			target.VersionsPath)
	}

	if updateCheckInterval == time.Duration(0) {
		return nil, fmt.Errorf(
			"UpdateCheckInterval value of '%v' is invalid",
			updateCheckInterval)
	}

	updater := Unattended{
		outputWriter:        os.Stdout,
		clientID:            clientID,
		target:              target,
		updateCheckInterval: updateCheckInterval,
		log:                 log,
	}

	return &updater, nil
}

// SetOutputWriter sets the writer to write the target's output to
func (updater *Unattended) SetOutputWriter(writer io.Writer) {
	updater.outputWriter = writer
}

// Run starts the target application and the update check loop.
//
// If any updates are found for targets in UpdateManifests they will be
// downloaded, applied and the target application restarted
func (updater *Unattended) Run() error {
	updater.log.WithField(
		"check_interval", updater.updateCheckInterval,
	).Info("Starting service with update checking enabled")
	time.AfterFunc(updater.updateCheckInterval, updater.handleUpdates)
	return updater.RunWithoutUpdate()
}

// RunWithoutUpdate starts the target application without checking for updates
func (updater *Unattended) RunWithoutUpdate() error {
	updater.command = exec.Command(
		filepath.Join(
			updater.target.VersionsPath,
			updater.target.LatestVersion(),
			updater.target.ApplicationName,
		),
		updater.target.ApplicationParameters...)
	// TODO: Do we need to set a process group on Linux?
	// updater.command.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	commandOutPipe, err := updater.command.StdoutPipe()
	if err != nil {
		return fmt.Errorf("Unable to start reading miner output")
	}
	commandErrPipe, err := updater.command.StderrPipe()
	if err != nil {
		return fmt.Errorf("Unable to start reading miner error output")
	}

	err = updater.command.Start()
	if err != nil {
		return err
	}

	// Keep copying the output from the process and send to the stream
	updater.waitGroup.Add(1)
	go func() {
		defer func() {
			updater.waitGroup.Done()
		}()
		io.Copy(updater.outputWriter, commandOutPipe)
		io.Copy(updater.outputWriter, commandErrPipe)
	}()

	updater.mutex.Lock()
	updater.commandCompleted = false
	updater.mutex.Unlock()

	err = updater.command.Wait()
	if err != nil {
		updater.log.Infof("Target completed: %s", err)
	}
	updater.mutex.Lock()
	updater.commandCompleted = true
	updater.mutex.Unlock()

	updater.waitGroup.Wait()
	return nil
}

// Stop the target application
func (updater *Unattended) Stop() error {
	updater.mutex.Lock()
	defer updater.mutex.Unlock()
	cmd := updater.command

	if cmd == nil {
		return nil
	}
	if cmd.Process == nil {
		return nil
	}

	updater.log.Infof("Stopping target, PID %d", cmd.Process.Pid)

	//
	// Simplified attempt at killing spree
	// os.Interrupt signal isn't available on Windows, so we're just
	// killing the processes
	//
	// This doesn't work reliably when unattended is used by a Windows service
	// that spawns a sub-unattended managed service. It needs some work.
	//
	if strings.ToLower(runtime.GOOS) == "linux" || strings.ToLower(runtime.GOOS) == "darwin" {
		// updater.log.Info("Releasing target process")
		// err := cmd.Process.Release()
		// if err != nil {
		// 	updater.log.Warningf("Target could not be released: %s", err)
		// 	// Return nil, there isn't much we can do now...
		// 	return nil
		// }
		updater.log.Info("Killing target process")
		err := cmd.Process.Kill()
		if err != nil {
			updater.log.Warningf("Target could not be killed: %s", err)
			// Return nil, there isn't much we can do now...
			return nil
		}

	} else if strings.ToLower(runtime.GOOS) == "windows" {
		updater.log.Info("Killing target process with taskkill")
		// Some processes just need to be force killed on Windows, many many
		// tests showed Windows not killing it when process.Kill is used.
		// This is especially true when this runs as a service
		_, err := exec.Command(
			"taskkill",
			"/F",   // Force
			"/PID", // by process ID
			fmt.Sprintf("%d", cmd.Process.Pid),
		).Output()
		if err != nil {
			updater.log.Warningf("Target could not be taskkilled: %s", err)
			// Return nil, there isn't much we can do now...
			return nil
		}
	}

	updater.log.Info("Target stopped")
	return nil
}

// Restart the target application
func (updater *Unattended) Restart() error {
	updater.log.Info("Restarting target")
	err := updater.Stop()
	if err != nil {
		// TODO: ROLLBACK
		return err
	}
	return updater.RunWithoutUpdate()
}

// handleUpdates runs at updateCheckInterval to check for and apply updates
func (updater *Unattended) handleUpdates() {

	updater.log.Debug("Checking for updates...")
	updated, err := updater.ApplyUpdates()
	if err != nil {
		updater.log.Warningf("Unable to check for updates: %s", err)
	}
	if updated {
		updater.log.WithField(
			"new_version", updater.target.LatestVersion(),
		).Info("Software updated")
		// Restart the application
		// TODO: Check if we are leaking a goroutine here
		go func() {
			err = updater.Restart()
			if err != nil {
				// TODO: ROLLBACK
				updater.log.Fatalf("Unable to restart target: %s", err)
			}
		}()
	} else {
		updater.log.Debug("No updates available")
	}
	time.AfterFunc(updater.updateCheckInterval, updater.handleUpdates)
}

// ApplyUpdates downloads and applies downloads if they are available
func (updater *Unattended) ApplyUpdates() (bool, error) {

	currentVersion := updater.target.LatestVersion()

	omahaManifests, err := updater.getAvailableUpdates()
	if err != nil {
		return false, fmt.Errorf("Unable to get updates: %s", err)
	}
	if len(omahaManifests) == 0 {
		return false, nil
	}

	updater.log.WithField(
		"updates", len(omahaManifests),
	).Debug("Updates found, download...")

	tempPath := filepath.Join(updater.target.VersionsPath, "tmp")
	//Remove/clean the temp directory
	err = os.RemoveAll(tempPath)
	if err != nil {
		updater.log.Warningf(
			"Unable to remove temp download path at '%s': %s",
			tempPath,
			err)
	}

	// And recreate/create the temp directory
	err = os.MkdirAll(tempPath, 0755)
	if err != nil {
		updater.log.Warningf(
			"Unable to create temp download path at '%s': %s",
			tempPath,
			err)
		return false, err
	}

	updater.log.WithField(
		"path", tempPath,
	).Debugf("Temp path set")

	for _, omahaManifest := range omahaManifests {
		downloadPath, err := updater.DownloadAndVerifyPackage(omahaManifest, tempPath)
		if err != nil {
			updater.log.WithFields(logrus.Fields{
				"package":         omahaManifest.Package.Name,
				"package_version": omahaManifest.Version,
				"reason":          err,
			}).Errorf("Unable to download package")

			continue
		}

		updater.log.WithFields(logrus.Fields{
			"package":         omahaManifest.Package.Name,
			"package_version": omahaManifest.Version,
		}).Debug("Downloaded package")

		// Get new version path
		newVersionPath := filepath.Join(updater.target.VersionsPath, omahaManifest.Version)
		updater.log.WithField(
			"path", newVersionPath,
		).Debugf("New version path set")

		// Clone the current version into new version
		// If no versions are currently installed, create the new path
		// Note: From this point on the new version folder might exist, in case of
		// rollback, remove this version
		currentVersionPath := filepath.Join(updater.target.VersionsPath, currentVersion)
		if _, err = os.Stat(currentVersionPath); err == nil {
			err = copy.Copy(
				filepath.Join(updater.target.VersionsPath, currentVersion),
				newVersionPath,
			)
			if err != nil {
				return false, updater.undoIncomplete(newVersionPath, err)
			}
		} else {
			// No current version exists, create the path
			err = os.MkdirAll(newVersionPath, 0755)
			if err != nil {
				return false, updater.undoIncomplete(newVersionPath, err)
			}
		}
		// Override files from package in new dir / apply update
		// Start with the gz part of the tar.gz file
		downloadedPackage, err := os.Open(downloadPath)
		if err != nil {
			return false, updater.undoIncomplete(newVersionPath, err)
		}
		gzReader, err := gzip.NewReader(downloadedPackage)
		if err != nil {
			return false, updater.undoIncomplete(newVersionPath, err)
		}

		tarReader := tar.NewReader(gzReader)
		// Go through all files in tar archive
		for {
			header, err := tarReader.Next()

			// No more
			if err == io.EOF {
				break
			}
			if err != nil {
				// Next!
				continue
			}

			// get the filename in the archive
			filename := header.Name
			destinationPath := filepath.Join(newVersionPath, filename)
			switch header.Typeflag {
			case tar.TypeDir:
				// Create directories if needed
				err := os.MkdirAll(destinationPath, header.FileInfo().Mode())
				if err != nil {
					return false, updater.undoIncomplete(newVersionPath, err)
				}
			case tar.TypeReg:
				// Create the new file in the destination path
				destinationFile, err := os.OpenFile(
					destinationPath,
					os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
					header.FileInfo().Mode())
				if err != nil {
					return false, updater.undoIncomplete(newVersionPath, err)
				}
				defer destinationFile.Close()

				written, err := io.Copy(destinationFile, tarReader)
				if err != nil {
					return false, updater.undoIncomplete(newVersionPath, err)
				}
				destinationFile.Close()

				if written != header.Size {
					return false, updater.undoIncomplete(newVersionPath, fmt.Errorf(
						"Written bytes differ from original file. Expected %d, wrote %d",
						header.Size,
						written))
				}
				updater.log.WithField(
					"path", destinationPath,
				).Debugf("Updated file")
			default:
				updater.log.Warningf("Unable to determine type, found: %c %s %s\n",
					header.Typeflag,
					"in file",
					filename,
				)
			}
		}
		err = gzReader.Close()
		if err != nil {
			updater.log.Warningf("Unable to close gz: %s", err)
		}

		err = downloadedPackage.Close()
		if err != nil {
			updater.log.Warningf("Unable to close package: %s", err)
		}

	}

	err = os.RemoveAll(tempPath)
	if err != nil {
		updater.log.Warningf("Unable to remove temp download path: %s", err)
	}

	return true, nil
}

// DownloadAndVerifyPackage downloads and verifies the package from the
// given manifest and returns the downloaded location
func (updater *Unattended) DownloadAndVerifyPackage(
	manifest omaha.Manifest,
	tempPath string) (string, error) {

	updater.log.WithField(
		"name", manifest.Package.Name,
	).Debugf("Downloading package")

	downloadPath := filepath.Join(tempPath, manifest.Package.Name)
	response, err := grab.Get(downloadPath, manifest.DownloadURL.Codebase)
	if err != nil {
		return "", err
	}

	hasher := sha256.New()
	downloadedFile, err := os.Open(response.Filename)
	if err != nil {
		return "", fmt.Errorf(
			"Unable to access: %s",
			err)
	}
	defer downloadedFile.Close()
	if _, err := io.Copy(hasher, downloadedFile); err != nil {
		return "", fmt.Errorf(
			"Could not be verified: %s",
			err)
	}

	if manifest.Package.SHA256Hash != hex.EncodeToString(hasher.Sum(nil)) {
		return "", fmt.Errorf("Failed verification")
	}

	return response.Filename, nil
}

// isUpdateAvailable checks if an update is available for the target and returns
// the available package if true
func (updater *Unattended) isUpdateAvailable() (bool, omaha.Manifest, error) {
	currentVersion := updater.target.LatestVersion()
	updater.log.WithFields(logrus.Fields{
		"app_id":          updater.target.AppID,
		"current_version": currentVersion,
		"update_endpoint": updater.target.UpdateEndpoint,
	}).Debug("Checking for update")

	// TODO: Convert unattended to proper semver
	// _, err := semver.Parse("1.0.0-dev")
	// if err != nil {
	// 	panic(err)
	// }
	//

	omahaRequest := omaha.Request{
		Protocol: 3,
		Application: omaha.App{
			Channel:  updater.target.UpdateChannel,
			ClientID: "1",
			ID:       updater.target.AppID,
			Version:  currentVersion,
			Event: omaha.Event{
				Type:   omaha.EventTypeUpdateCheck,
				Result: omaha.EventResultTypeStarted,
			},
		},
	}

	omahaBytes, err := xml.Marshal(omahaRequest)
	if err != nil {
		return false, omaha.Manifest{}, fmt.Errorf(
			"Unable to check for update, invalid request: %s",
			err)
	}

	response, err := http.Post(
		updater.target.UpdateEndpoint,
		"application/xml",
		bytes.NewReader(omahaBytes))
	if err != nil {
		return false, omaha.Manifest{}, fmt.Errorf(
			"Unable to check for update, received API error: %s",
			err)
	}

	if response.StatusCode != http.StatusOK {
		return false, omaha.Manifest{},
			fmt.Errorf(
				"Unable to check for update, received HTTP status code %d: %s",
				response.StatusCode,
				response.Status)
	}

	var omahaResponse omaha.Response
	err = xml.NewDecoder(response.Body).Decode(&omahaResponse)
	if err != nil {
		return false, omaha.Manifest{}, fmt.Errorf(
			"Unable to check for update, received invalid response: %s",
			err)
	}

	// Error getting update information
	if omahaResponse.Application.Status != "ok" {
		return false, omaha.Manifest{}, fmt.Errorf(
			"Received app status %s: %s",
			omahaResponse.Application.Status,
			omahaRequest.Application.Reason)
	}
	// No update is available
	if omahaResponse.Application.UpdateCheck.Status == "noupdate" {
		return false, omaha.Manifest{}, nil
	}
	if omahaResponse.Application.UpdateCheck.Status != "ok" {
		return false, omaha.Manifest{}, fmt.Errorf(
			"%s",
			omahaResponse.Application.UpdateCheck.Status)
	}
	return true, omahaResponse.Application.UpdateCheck.Manifest, nil
}

// getAvailableUpdates checks for all packages that have updates available
func (updater *Unattended) getAvailableUpdates() ([]omaha.Manifest, error) {
	var omahaManifests []omaha.Manifest

	hasUpdate, omahaManifest, err := updater.isUpdateAvailable()
	if err != nil {
		return omahaManifests, err
	}
	if hasUpdate {
		updater.log.WithFields(logrus.Fields{
			"app_id":            updater.target.AppID,
			"available_version": omahaManifest.Version,
		}).Debugf("Update available")
		omahaManifests = append(omahaManifests, omahaManifest)
	}

	return omahaManifests, nil
}

// GetLatestVersion returns the latest installed version
func (updater *Unattended) GetLatestVersion() string {
	return updater.target.LatestVersion()
}

// undoIncomplete removes an incomplete update
func (updater *Unattended) undoIncomplete(versionPath string, originalErr error) error {
	err := os.RemoveAll(versionPath)
	if err != nil {
		updater.log.Errorf("Unable to remove incomplete update: %s", err)
	}
	return originalErr
}

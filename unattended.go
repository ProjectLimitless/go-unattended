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
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/ProjectLimitless/go-unattended/omaha"
	"github.com/cavaliercoder/grab"
)

// Unattended implements the core functionality of the package. It takes
// ownership of running and updating a target application
type Unattended struct {
	// config of the Unattended update setup
	config Config
}

// New creates a new instance of the unattended updater
func New(config Config) (*Unattended, error) {

	updater := Unattended{
		config: config,
	}

	return &updater, nil
}

// Run starts the target application and the update check loop.
//
// If any updates are found for targets in UpdateManifests they will be
// downloaded, applied and the target application restarted
func (updater *Unattended) Run() error {
	fmt.Println("BLAH! RUN!")
	return nil
}

// ProcessUpdates checks, downloads and applies downloads if they are available
func (updater *Unattended) ProcessUpdates() error {

	omahaManifests, err := updater.getAvailableUpdates()
	if err != nil {
		return fmt.Errorf("Unable to get updates for all manifests: %s", err)
	}
	if len(omahaManifests) == 0 {
		fmt.Println("No updates are available")
		return nil
	}

	fmt.Printf("%d updates found.. download!\n", len(omahaManifests))

	tempPath := filepath.Join(path.Dir(updater.config.Target.Path), "../tmp")

	// Remove/clean the temp directory
	err = os.RemoveAll(tempPath)
	if err != nil {
		fmt.Printf("Unable to remove temp download path at '%s': %s\n",
			tempPath,
			err)
	}
	// And recreate/create the temp directory
	err = os.MkdirAll(tempPath, 0755)
	if err != nil {
		fmt.Printf("Unable to create temp download path at '%s': %s\n",
			tempPath,
			err)
	}

	fmt.Println("Downloading to", tempPath)

	for _, omahaManifest := range omahaManifests {
		downloadPath, err := updater.DownloadAndVerifyPackage(omahaManifest, tempPath)
		if err != nil {
			fmt.Printf("Unable to download package for '%s (%s)': %s\n",
				omahaManifest.Package.Name,
				omahaManifest.Version,
				err)
			continue
		}
		fmt.Printf("Downloaded package for '%s (%s)': %s\n",
			omahaManifest.Package.Name,
			omahaManifest.Version,
			downloadPath)

		// TODO: Ge latest version path
		// TODO: Clone the path into new version
		// TODO: Override files from package in new dir
		// TODO: Restart app
	}

	return nil
}

// DownloadAndVerifyPackage downloads and verifies the package from the
// given manifest and returns the downloaded location
func (updater *Unattended) DownloadAndVerifyPackage(
	manifest omaha.Manifest,
	tempPath string) (string, error) {

	fmt.Printf("Downloading package %s\n", manifest.Package.Name)

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
			manifest.Package.Name,
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

// IsUpdateAvailable checks if an update is available and returns the available
// package if true
func (updater *Unattended) IsUpdateAvailable(
	manifest UpdateManifest) (bool, omaha.Manifest, error) {

	currentVersionDirectory := path.Dir(updater.config.Target.Path)
	currentVersion := path.Base(currentVersionDirectory)

	fmt.Printf(
		"Checking updates for %s (%s) at %s\n",
		manifest.AppID,
		currentVersion,
		manifest.Endpoint)
	// TODO: Convert unattended to proper semver
	// _, err := semver.Parse("1.0.0-dev")
	// if err != nil {
	// 	panic(err)
	// }
	//

	omahaRequest := omaha.Request{
		Protocol: 3,
		Application: omaha.App{
			Channel:  "stable",
			ClientID: "1",
			ID:       manifest.AppID,
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
		manifest.Endpoint,
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
	fmt.Println("Checking for updates on all manifests")

	var omahaManifests []omaha.Manifest
	for _, updateManifest := range updater.config.UpdateManifests {
		hasUpdate, omahaManifest, err := updater.IsUpdateAvailable(updateManifest)
		if err != nil {
			return omahaManifests, err
		}
		if hasUpdate {
			fmt.Printf("Update available for %s (%s)\n",
				updateManifest.AppID,
				omahaManifest.Version)
			omahaManifests = append(omahaManifests, omahaManifest)
		}
	}

	return omahaManifests, nil
}

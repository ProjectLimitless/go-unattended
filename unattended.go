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
	"encoding/xml"
	"fmt"
	"net/http"
	"path"

	"github.com/ProjectLimitless/go-unattended/omaha"
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

// IsUpdateAvailable checks if an update is available and returns the available
// package if true
func (updater *Unattended) IsUpdateAvailable(
	manifest UpdateManifest) (bool, omaha.Manifest, error) {

	currentVersionDirectory := path.Dir(updater.config.Target.Path)
	currentVersion := path.Base(currentVersionDirectory)

	fmt.Printf(
		"Checking updates for %s (v%s) at %s\n",
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

	return true, omaha.Manifest{}, nil
}

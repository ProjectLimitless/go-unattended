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
	"fmt"

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

	fmt.Printf("Checking updates for %s at %s\n", manifest.AppID, manifest.Endpoint)

	return false, omaha.Manifest{}, nil
}

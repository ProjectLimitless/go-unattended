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
	"io/ioutil"
)

// Target defines the target application to be controlled and updated by
// Unattended
type Target struct {
	// AppID is the unique ID of the application to use in checking for updates
	AppID string
	// UpdateEndpoint is the Unattended server endpoint serving updates
	UpdateEndpoint string
	// UpdateChannel defines the update channel, can be 'stable', 'beta' or any
	// other value defined by the Unattended server
	UpdateChannel string
	// VersionsPath is the base path to where the versioned directories were
	// installed to
	VersionsPath string
	// ApplicationName is the name of the executable to run
	ApplicationName string
	// ApplicationParameters to use in executing the target
	ApplicationParameters []string
}

// LatestVersion returns the latest version installed of the target
func (target *Target) LatestVersion() string {
	// ReadDir lists the directories sorted, the last one will be
	// the latest version
	files, err := ioutil.ReadDir(target.VersionsPath)
	if err != nil {
		return ""
	}

	latestVersion := ""
	for _, f := range files {
		latestVersion = f.Name()
	}
	return latestVersion
}

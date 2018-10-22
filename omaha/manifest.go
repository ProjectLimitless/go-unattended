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

package omaha

import (
	"encoding/xml"
)

// Manifest of the update package
type Manifest struct {
	XMLName xml.Name `xml:"manifest"`
	// Version of the update
	Version string `xml:"version,attr"`
	// TraceID for the identification of this update
	TraceID string `xml:"trace,attr"`
	// DownloadURL is the location of the downloadable package
	DownloadURL URL `xml:"url"`
	// Package contains the validation information for the package
	// to be retrieved from DownloadUrl
	Package Package `xml:"package"`
}

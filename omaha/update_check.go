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

import "encoding/xml"

// UpdateCheck holds information for checking for an update
type UpdateCheck struct {
	XML xml.Name `xml:"updatecheck,omitempty"`
	// Status of the update check
	Status string `xml:"status,attr,omitempty"`
	// Manifest of the update package
	Manifest Manifest `xml:"manifest"`
}

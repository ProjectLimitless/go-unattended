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

// Package contains the validation information about an update
type Package struct {
	XMLName xml.Name `xml:"package"`
	// SHA256Hash of the download package
	SHA256Hash string `xml:"hash,attr,omitempty"`
	// Name of the download package
	Name string `xml:"name,attr,omitempty"`
	// SizeInBytes of the download package
	SizeInBytes uint64 `xml:"size,attr,omitempty"`
}

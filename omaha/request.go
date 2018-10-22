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

// Request contains a subset of the Omaha update request protocol
type Request struct {
	XMLName xml.Name `xml:"request"`
	// Protocol version of the request
	Protocol float32 `xml:"protocol,attr"`
	// Application information for the request
	Application App `xml:"app"`
}

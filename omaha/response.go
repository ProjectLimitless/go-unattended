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

// Response contains a subset of the Omaha request protocol
type Response struct {
	XMLName xml.Name `xml:"response"`
	// Protocol version fo the response
	Protocol float32 `xml:"protocol,attr"`
	// Application being responded on
	Application App `xml:"app"`
}

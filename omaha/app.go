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

// App defines the applicaiton information in the request
type App struct {
	XMLName xml.Name `xml:"app"`
	/// ID for the application
	ID string `xml:"appid,attr"`
	/// Status of the response
	Status string `xml:"status,attr,omitempty"`
	/// Version of the application in symver format
	Version string `xml:"version,attr"`
	/// Channel of the update, stable, beta, etc.
	Channel string `xml:"track,attr"`
	/// ClientID is the unique ID of the client
	ClientID string `xml:"bootid,attr"`
	/// Event being sent to the server
	Event Event `xml:"event"`
	/// Response for update events.
	UpdateCheck UpdateCheck `xml:"updatecheck,omitempty"`
	/// Reason for failure if status is not Ok
	Reason string `xml:"reason"`
}

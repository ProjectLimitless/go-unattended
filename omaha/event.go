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

const (
	// EventTypeUnknown is for unknown event types
	EventTypeUnknown string = "0"
	// EventTypeUpdateCheck is to check for update
	EventTypeUpdateCheck string = "1"
	// EventTypeDownload is for events regarding downloads
	EventTypeDownload string = "2"
	// EventTypeInstall is for events regarding installation
	EventTypeInstall string = "3"
	// EventTypeRollback is for events regarding rollback
	EventTypeRollback string = "4"
	// EventTypePing is for ping tests
	EventTypePing string = "800"
)

const (
	// EventResultTypeUnknown is for an unknown type
	EventResultTypeUnknown string = "0"
	// EventResultTypeNoUpdate is the result when no update is available
	EventResultTypeNoUpdate string = "1"
	// EventResultTypeAvailable the result when a new update is available
	EventResultTypeAvailable string = "2"
	// EventResultTypeSuccess is for operation success
	EventResultTypeSuccess string = "3"
	// EventResultTypeSuccessRestarted is for success and app restarted
	EventResultTypeSuccessRestarted string = "4"
	// EventResultTypeError is for operation failed
	EventResultTypeError string = "5"
	// EventResultTypeCancelled is for operation cancelled
	EventResultTypeCancelled string = "6"
	// EventResultTypeStarted is foro peration started
	EventResultTypeStarted string = "7"
)

// Event is the event being sent to the server
type Event struct {
	XMLName xml.Name `xml:"event"`
	// Type of event
	Type string `xml:"eventtype,attr"`
	// Result of event
	Result string `xml:"eventresult,attr"`
}

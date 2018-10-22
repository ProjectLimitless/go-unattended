package unattended

import "time"

// Config for Unattended
type Config struct {
	// ClientID is the unique id of this client to report on the user's dashboard
	ClientID string
	// Target application to run. control and update
	Target Target
	// UpdateCheckInterval defines the interval in which we
	// should check for updates
	UpdateCheckInterval time.Duration
	// UpdateManifests is a collection of updates to check at the given interval
	UpdateManifests []UpdateManifest
}

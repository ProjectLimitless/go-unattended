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
	// UpdateManifests is the manifests for updating packages with the target
	UpdateManifests []UpdateManifest
}

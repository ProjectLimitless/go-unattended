package unattended

// Target defines the target application to be controlled and updated by
// Unattended
type Target struct {
	// Path to the target application
	Path string
	// Parameters to use in executing the target
	Parameters []string
}

package unattended

// Unattended implements the core functionality of the package. It takes
// ownership of running and updating a target application
type Unattended struct {
	// Target application to control and update
	Target Target
}

// New creates a new instance of the unattended updater
func New() (*Unattended, error) {
	unattended := Unattended{}

	return &unattended, nil
}

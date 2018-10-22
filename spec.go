package unattended

import (
	"encoding/json"
	"io/ioutil"
)

// UpdateManifest contains the information about targets to check for updates
type UpdateManifest struct {
	// AppID is the unique identifier of the target to check updates for
	AppID string
	// Endpoint specifies the URL of the Unattended server API
	Endpoint string
}

// NewUpdateManifestManifestFromFile creates a new
// update manifest from a JSON file
func NewUpdateManifestManifestFromFile(path string) (UpdateManifest, error) {
	var manifest UpdateManifest

	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return manifest, err
	}

	err = json.Unmarshal(fileBytes, &manifest)
	return manifest, err
}

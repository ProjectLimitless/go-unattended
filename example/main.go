// Package main implements a simple sample of the unattended package
package main

import (
	"fmt"
	"time"

	unattended "github.com/ProjectLimitless/go-unattended"
)

func main() {
	fmt.Println("Simple update sample")

	config := unattended.Config{
		ClientID: "TEST001",
		Target: unattended.Target{
			Path: "ping",
		},
		UpdateCheckInterval: time.Minute,
		UpdateManifests: []unattended.UpdateManifest{
			{
				AppID:    "testapp",
				Endpoint: "http://unattended.local/api",
			},
		},
	}

	updater, err := unattended.New(config)
	if err != nil {
		panic(err)
	}

	canUpdate, updatePackage, err := updater.IsUpdateAvailable(
		unattended.UpdateManifest{
			AppID:    "testapp",
			Endpoint: "http://unattended.local/api",
		},
	)

	if err != nil {
		panic(err)
	}
	if canUpdate {
		fmt.Println("Update available!")
		fmt.Println(updatePackage.Package.Name, updatePackage.Version)
	}
}

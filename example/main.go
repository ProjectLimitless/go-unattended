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
		ClientID: "GoTEST001",
		Target: unattended.Target{
			Path: "/home/donovan/Development/Go/code/src/github.com/ProjectLimitless/go-unattended/example/apptoupdate/1.0.0.0/apptoupdate",
		},
		UpdateCheckInterval: time.Minute,
		UpdateManifests: []unattended.UpdateManifest{
			unattended.UpdateManifest{
				AppID:    "apptoupdate",
				Endpoint: "http://unattended-old.local",
			},
		},
	}

	updater, err := unattended.New(config)
	if err != nil {
		panic(err)
	}

	canUpdate, updatePackage, err := updater.IsUpdateAvailable(
		unattended.UpdateManifest{
			AppID:    "apptoupdate",
			Endpoint: "http://unattended-old.local",
		},
	)

	if err != nil {
		panic(err)
	}
	if canUpdate {
		fmt.Println("Update available!")
		fmt.Println(updatePackage.Package.Name, updatePackage.Version)
	} else {
		fmt.Println("No update")
	}

	fmt.Println("Process all updates")
	err = updater.ProcessUpdates()
	if err != nil {
		panic(err)
	}
}

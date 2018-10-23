// Package main implements a simple sample of the unattended package
package main

import (
	"fmt"
	"os"
	"time"

	unattended "github.com/ProjectLimitless/go-unattended"
	"github.com/sirupsen/logrus"
)

func main() {

	// TODO: Check if we are in debug mode?
	debug := true

	// Setup the logging, by default we log to stdout
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "Jan 02 15:04:05",
	})
	logrus.SetLevel(logrus.InfoLevel)
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	logrus.SetOutput(os.Stdout)
	log := logrus.WithFields(logrus.Fields{
		"service": "unattended-test",
	})
	log.Info("Setting up Unattended updates")

	updater, err := unattended.New(
		"GoTEST001", // clientID
		unattended.Target{ // target
			VersionsPath:    "./apptoupdate",
			AppID:           "apptoupdate",
			UpdateEndpoint:  "http://unattended-old.local",
			UpdateChannel:   "stable",
			ApplicationName: "apptoupdate",
		},
		time.Second*5, // UpdateCheckInterval
		log,
	)
	if err != nil {
		panic(err)
	}

	go func() {

		err = updater.Run()
		if err != nil {
			// TODO: LOG!
			fmt.Println("Run error: ", err)
		}
	}()

	time.Sleep(time.Second * 8)
	fmt.Println("Stopping")
	err = updater.Stop()
	if err != nil {
		panic(err)
	}
}

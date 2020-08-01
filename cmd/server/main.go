package main

import (
	"os"
	"our-expenses-server/container"

	"github.com/sirupsen/logrus"
)

func main() {
	mongoDatabase, mongoError := container.InitDatabase()
	if mongoError != nil {
		logrus.Fatalf("Failed establish MongoDB connection: '%s'", mongoError)
	}

	server, serverError := container.CreateServer(mongoDatabase)
	if serverError != nil {
		logrus.Fatalf("Failed to start web server: '%s'", serverError)
	}
	server.Start()

	os.Exit(0)
}

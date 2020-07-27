package main

import (
	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	mongoDatabase, mongoError := initDatabase()
	if mongoError != nil {
		logrus.Fatalf("Failed establish MongoDB connection: '%s'", mongoError)
	}

	server, serverError := createServer(mongoDatabase)
	if serverError != nil {
		logrus.Fatalf("Failed to start web server: '%s'", serverError)
	}
	server.Start()

	os.Exit(0)
}

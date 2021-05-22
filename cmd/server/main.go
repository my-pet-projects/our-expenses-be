package main

import (
	"os"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/container"

	// "github.com/lightstep/otel-launcher-go/launcher"
	"github.com/lightstep/otel-launcher-go/launcher"
	"github.com/sirupsen/logrus"
)

func main() {
	// TODO: initialize logger and config here and pass as a dependency everywhere

	otel := initLightstepTracing()
	defer otel.Shutdown()

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

func initLightstepTracing() launcher.Launcher {
	launcher := launcher.ConfigureOpentelemetry(
		launcher.WithLogLevel("debug"),
		launcher.WithServiceName("our-expenses-server"),
		launcher.WithAccessToken("vAb0nB0w2ib4fMEa6VdZ0X47ZpoUJCjmu1xvSB8mT9JIQe3oBdvOpB5hmhKa3M5RddRqwfbHwsJZX5LlvK2XiAKNP/eO9KSyFB+wJ6bM"),
		// launcher.WithLogger(log),
	)
	logrus.Info("Initialized Lightstep OpenTelemetry launcher")
	return launcher
}

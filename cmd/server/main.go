package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/app"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/ports"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/server"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	app, appErr := app.NewApplication(ctx, cancel)
	if appErr != nil {
		log.Printf("Failed to instantiate application: %v", appErr)
		os.Exit(1)
	}

	signChan := make(chan os.Signal, 1)
	signal.Notify(signChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		signal := <-signChan
		app.Logger.Warnf(ctx, "Received interuption signal: %s", signal)
		cancel()
	}()

	server := server.NewServer(app.Logger, app.Config.Server,
		func(router *echo.Echo) {
			ports.RegisterHandlersWithBaseURL(router, ports.NewHTTPServer(app), "api")
		})
	if err := server.Start(ctx); err != nil {
		app.Logger.Error(ctx, "Failed to start HTTP server", err)
		os.Exit(1)
	}

	app.Tracer.Shutdown()
	app.Logger.Info(ctx, "Application shutdown")

	defer os.Exit(0)
}

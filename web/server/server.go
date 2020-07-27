package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"our-expenses-server/config"
	"our-expenses-server/logger"
	"our-expenses-server/web/api"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
)

// Server is a wrapper around an HTTP server.
type Server struct {
	server *http.Server
	logger *logger.AppLogger
}

// ProvideServer returns server instance with configured routes.
func ProvideServer(config *config.Config, logger *logger.AppLogger, router *api.Router) (*Server, error) {
	routers := router.InitializeRoutes()
	loggedRouter := handlers.LoggingHandler(logger.Writer(), routers)

	server := &http.Server{
		Addr:         ":" + fmt.Sprint(config.Port),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
		Handler:      loggedRouter,
	}
	return &Server{server: server, logger: logger}, nil
}

// Start spin-ups the web server.
func (srv Server) Start() {
	go func() {
		srv.logger.Info("Starting web server ...", logger.Fields{})
		if serverError := srv.server.ListenAndServe(); serverError != nil {
			srv.logger.Fatal("Could not start the server", serverError, logger.Fields{})
		}
		srv.logger.Info(fmt.Sprintf("Server is up and running on %v", srv.server.Addr), logger.Fields{})
	}()

	gracefulStop := make(chan os.Signal, 1)
	signal.Notify(gracefulStop, os.Interrupt, syscall.SIGTERM)

	killSignal := <-gracefulStop

	srv.logger.Info(fmt.Sprintf("Server is shutting down, reason: %v", killSignal.String()), logger.Fields{})

	serverCtx, serverCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer serverCancel()

	srv.server.SetKeepAlivesEnabled(false)
	if serverError := srv.server.Shutdown(serverCtx); serverError != nil {
		srv.logger.Fatal("Could not gracefully shutdown the server", serverError, logger.Fields{})
	}

	srv.logger.Info("Server stopped", logger.Fields{})
}

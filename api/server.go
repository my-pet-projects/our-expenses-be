package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"our-expenses-server/api/middleware"
	"our-expenses-server/api/router"
	"our-expenses-server/config"
	"our-expenses-server/logger"
	"syscall"
	"time"
)

// Server is a wrapper around an HTTP server.
type Server struct {
	server *http.Server
	logger *logger.AppLogger
}

// ProvideServer returns server instance with configured routes.
func ProvideServer(config *config.Config, logger *logger.AppLogger, router *router.Router) (*Server, error) {
	loggingMiddleware := middleware.LoggingMiddleware(logger)
	correlationIDMiddleware := middleware.CorrelationMiddleware(logger)

	routers := router.InitializeRoutes()
	loggedRouter := loggingMiddleware(routers)
	correlationIDRouter := correlationIDMiddleware(loggedRouter)
	corsedRouter := middleware.CorsMiddleware(correlationIDRouter)

	server := &http.Server{
		Addr:         "localhost:" + fmt.Sprint(config.Port),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
		Handler:      corsedRouter,
	}
	return &Server{server: server, logger: logger}, nil
}

//https://github.com/obiwandsilva/go-bestflight/blob/ffa0ec357d16d754217b77fafce404f48f7524cc/application/application.go

// Start spin-ups the web server.
func (srv Server) Start() {
	ctx := context.Background()
	go func() {
		if serverError := srv.server.ListenAndServe(); serverError != nil {
			srv.logger.Fatal("Could not start the server", serverError, logger.FieldsSet{})
		}
	}()

	srv.logger.Infof(ctx, "Server is up and running on %s", srv.server.Addr)

	gracefulStop := make(chan os.Signal, 1)
	signal.Notify(gracefulStop, os.Interrupt, syscall.SIGTERM)

	killSignal := <-gracefulStop

	srv.logger.Infof(ctx, "Server is shutting down, reason: %s", killSignal.String())

	serverCtx, serverCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer serverCancel()

	srv.server.SetKeepAlivesEnabled(false)
	if serverError := srv.server.Shutdown(serverCtx); serverError != nil {
		srv.logger.Fatal("Could not gracefully shutdown the server", serverError, logger.FieldsSet{})
	}

	srv.logger.Info(ctx, "Server stopped")
}

// func GracefullShutdown(quitChan chan os.Signal) {
// 	go func() {
// 		log.Println("gracefull shutdown enabled")

// 		oscall := <-quitChan

// 		log.Printf("system call:%+v", oscall)
// 		log.Println("shutting server down...")

// 		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 		defer cancel()

// 		server.SetKeepAlivesEnabled(false)

// 		if err := server.Shutdown(ctx); err != nil {
// 			log.Fatal("erro when shuttingdown the server:", err)
// 		}

// 		log.Println("finished")
// 		log.Println("press ctrl + D to exit program completely")
// 		fmt.Println("\nserver shutted down. Press ctrl + D to exit program completely")
// 	}()
// }

// server := &http.Server{
// 	Addr:         listenAddr,
// 	Handler:      tracing(nextRequestID)(logging(logger)(router)),
// 	ErrorLog:     logger,
// 	ReadTimeout:  5 * time.Second,
// 	WriteTimeout: 10 * time.Second,
// 	IdleTimeout:  15 * time.Second,
// }

// done := make(chan bool)
// quit := make(chan os.Signal, 1)
// signal.Notify(quit, os.Interrupt)

// go func() {
// 	<-quit
// 	logger.Println("Server is shutting down...")
// 	atomic.StoreInt32(&healthy, 0)

// 	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
// 	defer cancel()

// 	server.SetKeepAlivesEnabled(false)
// 	if err := server.Shutdown(ctx); err != nil {
// 		logger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
// 	}
// 	close(done)
// }()

// logger.Println("Server is ready to handle requests at", listenAddr)
// atomic.StoreInt32(&healthy, 1)
// if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
// 	logger.Fatalf("Could not listen on %s: %v\n", listenAddr, err)
// }

// <-done
// logger.Println("Server stopped")

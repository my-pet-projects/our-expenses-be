package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/config"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

// Server provides an HTTP server.
type Server struct {
	httpServer *echo.Echo
	logger     logger.LogInterface
	config     config.Server
}

// NewServer returns server instance with configured routes.
func NewServer(logger logger.LogInterface, config config.Server, registerHandlers func(e *echo.Echo)) *Server {
	ctx, cancel := context.WithCancel(context.Background())

	e := echo.New()
	registerHandlers(e)
	registerMiddleware(config, e)

	e.HideBanner = true
	e.HidePort = true
	e.Server = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
		ReadTimeout:  time.Duration(config.Timeout.Read) * time.Second,
		WriteTimeout: time.Duration(config.Timeout.Write) * time.Second,
		IdleTimeout:  time.Duration(config.Timeout.Idle) * time.Second,
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
	}
	e.Server.RegisterOnShutdown(cancel)

	return &Server{
		httpServer: e,
		logger:     logger,
		config:     config,
	}
}

func registerMiddleware(config config.Server, e *echo.Echo) {
	e.Use(otelecho.Middleware(config.Name))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPost},
	}))
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: func(e echo.Context) bool {
			return strings.Contains(e.Path(), "health")
		},
		Format: `{"time":"${time_rfc3339_nano}","remote_ip":"${remote_ip}",` +
			`"host":"${host}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}",` +
			`"status":${status},"error":"${error}","latency":"${latency_human}",` +
			`"bytes_in":${bytes_in},"bytes_out":${bytes_out},"severity":"info"}` + "\n",
		CustomTimeFormat: "2006-01-02 15:04:05.00000",
	}))
	e.Use(middleware.Recover())
}

// Start spin-ups the web server with graceful shutdown.
func (srv Server) Start(ctx context.Context) error {
	srv.logger.Info(ctx, "Starting server...")
	srv.logger.Infof(ctx, "Server is up and running on %s", srv.httpServer.Server.Addr)

	errChan := make(chan error, 1)

	go func() {
		errChan <- srv.shutdown(ctx)
	}()

	if err := srv.httpServer.Start(fmt.Sprintf("%s:%d", srv.config.Host, srv.config.Port)); err != http.ErrServerClosed {
		return errors.Wrap(err, "failed to start http server")
	}

	return <-errChan
}

func (srv Server) shutdown(ctx context.Context) error {
	<-ctx.Done()

	srv.httpServer.Server.SetKeepAlivesEnabled(false)

	timeout := time.Duration(srv.config.Timeout.Shutdown) * time.Second
	srv.logger.Infof(ctx, "Shutting server with %s timeout\n", timeout)
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), timeout)
	defer shutdownCancel()

	err := srv.httpServer.Shutdown(shutdownCtx)
	if err == context.DeadlineExceeded {
		srv.logger.Warn(ctx, "Some open connections were interrupted during shutdown timeout")
		err = nil
	} else if err != nil {
		return errors.Wrap(err, "gracefully shutdown the server")
	}

	srv.logger.Info(ctx, "Server gracefully stopped")
	return err
}

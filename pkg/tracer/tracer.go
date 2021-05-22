package tracer

import (
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/config"
	"github.com/lightstep/otel-launcher-go/launcher"
)

// Tracer is a wrapper around Lightstep tracer.
type Tracer struct {
	launcher launcher.Launcher
}

// TraceInterface defines tracer contract.
type TraceInterface interface {
	Shutdown()
}

// NewTracer instantiates Lightstep tracer.
func NewTracer(config config.Telemetry) *Tracer {
	launcher := launcher.ConfigureOpentelemetry(
		launcher.WithLogLevel(config.Level),
		launcher.WithServiceName(config.ServiceName),
		launcher.WithAccessToken(config.AccessToken),
	)
	return &Tracer{
		launcher: launcher,
	}
}

// Shutdown shutsdown tracer.
func (t Tracer) Shutdown() {
	t.launcher.Shutdown()
}

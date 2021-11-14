package tracer

import (
	"github.com/lightstep/otel-launcher-go/launcher"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/config"
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
		launcher.WithServiceName(config.ServiceName),
		launcher.WithAccessToken(config.AccessToken),
		launcher.WithLogLevel(config.Level),
	)

	return &Tracer{
		launcher: launcher,
	}
}

// Shutdown shutsdowns the tracer.
func (t Tracer) Shutdown() {
	t.launcher.Shutdown()
}

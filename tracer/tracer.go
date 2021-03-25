package tracer

import (
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// InitJaeger creates new trace exporter:
// - sends traces using Jaeger UDP proto
// - sample all traces (to use with Loki)
func InitJaeger(name, host, port string) func() {
	// create and install Jaeger export pipeline.
	flush, err := jaeger.InstallNewPipeline(
		jaeger.WithAgentEndpoint(fmt.Sprintf("%s:%s", host, port)),
		jaeger.Process{
			ServiceName: name,
		},
		jaeger.WithSDK(&sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
	)
	if err != nil {
		panic(err)
	}

	// propagate trace context
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}),
	)

	return flush
}

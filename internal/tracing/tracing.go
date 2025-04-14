package tracing

import (
	"fmt"
	"io"
	"log"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
)

type TracerConfig struct {
	ServiceName string
	Host        string
	Port        int
}

func InitTracer(tConfig TracerConfig) (opentracing.Tracer, io.Closer) {
	cfg := config.Configuration{
		ServiceName: tConfig.ServiceName,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: fmt.Sprintf("%s:%d", tConfig.Host, tConfig.Port),
		},
	}

	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		log.Fatalf("Failed to init tracer: %v", err)
	}
	opentracing.SetGlobalTracer(tracer)

	return tracer, closer
}

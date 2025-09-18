package tracer

import (
	"fmt"
	"io"

	"github.com/BagusAK95/amarta_test/internal/config"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

func Init(cfg config.JaegerConfig) (opentracing.Tracer, io.Closer) {
	cfgSource := &jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		},
	}

	jcfg, err := cfgSource.FromEnv()
	if err != nil {
		panic(fmt.Sprintf("cannot parse jaeger env: %v", err))
	}

	jcfg.ServiceName = cfg.ServiceName

	tracer, closer, err := jcfg.NewTracer()
	if err != nil {
		panic(fmt.Sprintf("cannot create new tracer: %v", err))
	}

	return tracer, closer
}

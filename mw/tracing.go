package mw

import (
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
	"io"
	"log"
	"net/http"
)

var tracer opentracing.Tracer

func setTracer(serviceName string) io.Closer {
	cfg := jaegercfg.Configuration{
		ServiceName: serviceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
		},
	}
	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory
	var closer io.Closer
	tracer, closer, _ = cfg.NewTracer(
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
	)
	opentracing.SetGlobalTracer(tracer)
	return closer
}

func SetServerSpan(spanName string, address string, pattern string) {
	closer := setTracer(spanName)
	defer closer.Close()
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		serverSpan := tracer.StartSpan(spanName, ext.RPCServerOption(spanCtx))
		defer serverSpan.Finish()
	})
	log.Fatal(http.ListenAndServe(address, nil))
}

func SetClientSpan(spanName string, url string, requestType string) *http.Response {
	closer := setTracer(spanName)
	defer closer.Close()

	clientSpan := tracer.StartSpan(spanName)
	defer clientSpan.Finish()
	req, _ := http.NewRequest(requestType, url, nil)
	ext.SpanKindRPCClient.Set(clientSpan)
	ext.HTTPUrl.Set(clientSpan, url)
	ext.HTTPMethod.Set(clientSpan, requestType)
	tracer.Inject(clientSpan.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	resp, _ := http.DefaultClient.Do(req)
	return resp
}

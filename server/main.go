package main

import (
	middleware "github.com/MagalixTechnologies/core/middleware"
	goahttp "goa.design/goa/v3/http"
	"jaeger-tracing/mw"
	"net/http"
)

func main() {
	var mux goahttp.Muxer
	{
		mux = goahttp.NewMuxer()
	}

	var handler http.Handler = mux
	{
		handler = mw.SetServerSpan("server")(handler)
		handler = middleware.Log(middleware.InfoLevel)(handler)
	}

	srv := &http.Server{Addr: "localhost:8082", Handler: handler}
	srv.ListenAndServe()
}

package main

import "jaeger-tracing/mw"

func main() {
	resp := mw.SetClientSpan("client", "http://localhost:8082/publish", "GET")
	print(resp.StatusCode)
}

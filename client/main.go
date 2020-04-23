package main

import "github.com/akhidrb/jaeger-tracer/mw"

func main() {
	resp := mw.SetClientSpan("client", "http://localhost:8082/publish", "GET")
	print(resp.StatusCode)
}

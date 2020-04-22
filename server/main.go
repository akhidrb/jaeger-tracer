package main

import (
	"jaeger-tracing/mw"
)
func main() {
	mw.SetServerSpan("server", "localhost:8082", "/publish")
}

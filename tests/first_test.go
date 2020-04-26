package tests

// Sample test to show how to write a basic unit test.
import (
	middleware "github.com/MagalixTechnologies/core/middleware"
	"github.com/akhidrb/jaeger-tracer/mw"
	"github.com/go-chi/chi"
	"net/http"
	"testing"
)

const checkMark = "\u2713"
const ballotX = "\u2717"

func TestTracing(t *testing.T) {
	go func() {
		r := chi.NewRouter()
		r.Use(mw.SetServerSpan("server"))
		r.Use(middleware.Log(middleware.InfoLevel))

		r.Get("/publish", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		http.ListenAndServe("localhost:8082", r)
	}()
	resp := mw.SetClientSpan("client", "http://localhost:8082/publish", "GET")
	statusCode := 200
	if resp.StatusCode == statusCode {
		t.Logf("\t\tShould receive a \"%d\" status. %v",
			statusCode, checkMark)
	} else {
		t.Errorf("\t\tShould receive a \"%d\" status. %v %v",
			statusCode, ballotX, resp.StatusCode)
	}
}

package exporter

import (
	"net/http"

	"neverending.dev/weather/ecowitt"
)

func Healthcheck(w http.ResponseWriter, r *http.Request) {
	switch ecowitt.WS.Status {
	case ecowitt.Ready:
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	case ecowitt.NotReady:
		w.WriteHeader(404)
		w.Write([]byte("NOT READY"))
	default:
		w.WriteHeader(503)
		w.Write([]byte("Unavailable"))
	}
}

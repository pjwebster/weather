package main

import (
	"log"
	"net/http"

	"neverending.dev/weather/ecowitt"
	"neverending.dev/weather/exporter"
)

func main() {

	http.Handle("/", http.FileServer(http.Dir("./dist")))
	http.HandleFunc("/healthz", exporter.Healthcheck)
	http.HandleFunc("/metrics", exporter.Serve)
	http.HandleFunc("/weather", ecowitt.ReportHandler)

	log.Fatal(http.ListenAndServe(":8090", nil))
}

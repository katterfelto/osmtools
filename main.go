package main

import (
	"fmt"
	"log/slog"
	"net/http"
)

func main() {
	log := slog.Default()

	var port int = 8000
	log.Info("Initializing server", "port", port)

	mux := http.NewServeMux()

	// Create the route handlers and associate them with the corresponding paths
	mux.HandleFunc("GET /", getHome)
	mux.HandleFunc("POST /coop30day-osm", postCoop30DayToOsmHandler)
	mux.HandleFunc("POST /coopcustom-osm", postCoopCustomToOsmHandler)

	log.Info("Server starting")
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux); err != nil {
		log.Error("Server failed to start:", "error", err)
	}
	log.Info("Server stopped")
}

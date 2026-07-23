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
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), securityHeaders(mux)); err != nil {
		log.Error("Server failed to start:", "error", err)
	}
	log.Info("Server stopped")
}

// securityHeaders wraps a handler and sets common security-related HTTP
// response headers on every request.
func securityHeaders(next http.Handler) http.Handler {
	const csp = "default-src 'self'; " +
		"style-src 'self' https://cdn.jsdelivr.net; " +
		"script-src 'self' https://unpkg.com 'sha256-FruV5A/1bRz0wcIwjTzrQ4MrCrBRSaobjSzOBRq4mSY='; " +
		"img-src 'self' data:; " +
		"font-src 'self'; " +
		"connect-src 'self'; " +
		"worker-src 'self' blob:; " +
		"form-action 'self'; " +
		"base-uri 'self'; " +
		"object-src 'none'; " +
		"frame-ancestors 'self'"

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("Content-Security-Policy", csp)
		h.Set("X-Frame-Options", "DENY")
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		h.Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		h.Set("Permissions-Policy", "camera=(), microphone=(), geolocation=(), payment=(), usb=()")

		next.ServeHTTP(w, r)
	})
}

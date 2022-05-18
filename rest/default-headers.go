package rest

import "net/http"

func defaultHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src: 'none'; style-src 'unsafe-inline';")
}

package gaugin_middleware

import "net/http"

const (
	htmlContentType = "text/html"
	defaultCSP      = "default-src 'self'; " +
		"script-src 'none'; " +
		"object-src 'none'; " +
		"style-src 'unsafe-inline';"
)

func DefaultHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", htmlContentType)
	w.Header().Set("Content-Security-Policy", defaultCSP)
}

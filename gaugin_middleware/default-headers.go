package gaugin_middleware

import "net/http"

const (
	htmlContentType = "text/html"
	defaultCSP      = "default-src 'self'; " +
		"script-src " +
		//script-product-details.gohtml
		"'sha256-89hiWWr+KJhZ6UNrcqHjAG0875M2T59Rm6nP6YF3/F8=' " +
		//script-iframe.gohtml
		"'sha256-vEdzDTUjeRFG21L/pW+qldt1k+gnTSWl4v2E16iqJPc=' " +
		"none; " +
		"object-src 'none'; " +
		"img-src 'self' data:; " +
		"style-src 'unsafe-inline';"
)

func DefaultHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", htmlContentType)
	w.Header().Set("Content-Security-Policy", defaultCSP)
}

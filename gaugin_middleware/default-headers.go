package gaugin_middleware

import "net/http"

const (
	htmlContentType = "text/html"
	defaultCSP      = "default-src 'self'; " +
		"script-src " +
		//script-product-details.gohtml
		"'sha256-00g6CyN+3M94jXLajJMHpAH2oPJlfW7tn+iQP49ERGs=' " +
		//script-iframe.gohtml
		"'sha256-vEdzDTUjeRFG21L/pW+qldt1k+gnTSWl4v2E16iqJPc=' " +
		"none; " +
		"object-src 'none'; " +
		"style-src 'unsafe-inline';"
)

func DefaultHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", htmlContentType)
	w.Header().Set("Content-Security-Policy", defaultCSP)
}

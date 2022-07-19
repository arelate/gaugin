package gaugin_middleware

import "net/http"

const (
	htmlContentType = "text/html"
	defaultCSP      = "default-src 'self'; " +
		"script-src " +
		//script-product-details.gohtml
		"'sha256-EoiesIg5jhsIaHn7PSaZ/oT9Yi0MCUx9WzALOyH9HkE=' " +
		//script-iframe.gohtml
		"'sha256-vEdzDTUjeRFG21L/pW+qldt1k+gnTSWl4v2E16iqJPc=' " +
		//script-missing-images.gohtml
		"'sha256-vAj29f8vCI+Qq0UyQqV6YB/RcBCG4wjrOhqmemg5uoU=' " +
		"none; " +
		"object-src 'none'; " +
		"img-src 'self' data:; " +
		"style-src 'unsafe-inline';"
)

func DefaultHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", htmlContentType)
	w.Header().Set("Content-Security-Policy", defaultCSP)
}

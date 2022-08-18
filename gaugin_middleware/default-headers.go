package gaugin_middleware

import (
	"fmt"
	"net/http"
	"strings"
)

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
		//script-image-fade-in.gothml
		"'sha256-0o6kCVC1+8tieHTMZdKurbNC2fDP/bQEO5AAmFyRqBo=' " +
		"'unsafe-inline'; " +
		"object-src 'none'; " +
		"img-src 'self' data:; " +
		"style-src 'unsafe-inline';"
)

func DefaultHeaders(timing map[string]int64, w http.ResponseWriter) {
	w.Header().Set("Content-Type", htmlContentType)
	w.Header().Set("Content-Security-Policy", defaultCSP)
	if timing != nil {
		stb := strings.Builder{}
		for topic, dur := range timing {
			stb.WriteString(fmt.Sprintf("%s;dur=%d,", topic, dur))
		}
		if stb.Len() > 0 {
			w.Header().Set("Server-Timing", stb.String())
		}
	}
}

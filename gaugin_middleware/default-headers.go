package gaugin_middleware

import (
	"github.com/boggydigital/stencil"
	"net/http"
	"strings"
)

const (
	htmlContentType = "text/html"
	defaultCSP      = "default-src 'self'; " +
		"object-src 'none'; " +
		"img-src 'self' data:; " +
		"style-src 'unsafe-inline';"
)

func DefaultHeaders(st *ServerTimings, w http.ResponseWriter) {
	w.Header().Set("Content-Type", htmlContentType)
	w.Header().Set("Cache-Control", "max-age=3600")
	stencilCSP := defaultCSP + "script-src " + strings.Join(stencil.ScriptHashes, " ") +
		//script-image-fade-in.gothml
		" 'sha256-0o6kCVC1+8tieHTMZdKurbNC2fDP/bQEO5AAmFyRqBo=' "

	w.Header().Set("Content-Security-Policy", stencilCSP)
	if st != nil {
		w.Header().Add("Server-Timing", st.String())
	}
}

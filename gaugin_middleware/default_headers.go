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

func DefaultHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", htmlContentType)
	stencilCSP := defaultCSP + "script-src " + strings.Join(stencil.ScriptHashes, " ")
	w.Header().Set("Content-Security-Policy", stencilCSP)
}

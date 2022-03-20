package rest

import (
	"github.com/boggydigital/nod"
	"net/http"
)

func HandleFuncs() {

	patternHandlers := map[string]http.Handler{
		// start at the account
		"/": http.RedirectHandler("/news", http.StatusPermanentRedirect),
		// current endpoints
		"/news":        Gzip(nod.RequestLog(http.HandlerFunc(GetNews))),
		"/downloads":   Gzip(nod.RequestLog(http.HandlerFunc(GetDownloads))),
		"/all":         Gzip(nod.RequestLog(http.HandlerFunc(GetAll))),
		"/product":     Gzip(nod.RequestLog(http.HandlerFunc(GetProduct))),
		"/find":        Gzip(nod.RequestLog(http.HandlerFunc(GetFind))),
		"/images":      nod.RequestLog(http.HandlerFunc(GetImages)),
		"/videos":      nod.RequestLog(http.HandlerFunc(GetVideos)),
		"/files":       basicHttpAuth(nod.RequestLog(http.HandlerFunc(GetFiles))),
		"/local-file/": basicHttpAuth(nod.RequestLog(http.HandlerFunc(GetLocalFile))),
		"/favicon.ico": http.HandlerFunc(http.NotFound),
	}

	for p, h := range patternHandlers {
		http.HandleFunc(p, h.ServeHTTP)
	}
}

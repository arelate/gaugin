package api

import (
	"github.com/boggydigital/nod"
	"net/http"
)

func HandleFuncs() {

	patternHandlers := map[string]http.Handler{
		// legacy endpoints
		"/account": http.RedirectHandler("/downloadable", http.StatusPermanentRedirect),
		"/store":   http.RedirectHandler("/all", http.StatusPermanentRedirect),
		// start at the account
		"/": http.RedirectHandler("/downloadable", http.StatusPermanentRedirect),
		//
		"/downloadable": Gzip(memCache(nod.RequestLog(http.HandlerFunc(GetDownloadable)))),
		"/all":          Gzip(memCache(nod.RequestLog(http.HandlerFunc(GetAll)))),
		"/product":      Gzip(nod.RequestLog(http.HandlerFunc(GetProduct))),
		"/search":       Gzip(nod.RequestLog(http.HandlerFunc(GetSearch))),
		"/images":       nod.RequestLog(http.HandlerFunc(GetImages)),
		"/videos":       nod.RequestLog(http.HandlerFunc(GetVideos)),
		"/files":        basicHttpAuth(nod.RequestLog(http.HandlerFunc(GetFiles))),
		"/local-file/":  basicHttpAuth(nod.RequestLog(http.HandlerFunc(GetLocalFile))),
		"/favicon.ico":  http.HandlerFunc(http.NotFound),
	}

	for p, h := range patternHandlers {
		http.HandleFunc(p, h.ServeHTTP)
	}
}

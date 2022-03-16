package api

import (
	"github.com/boggydigital/nod"
	"net/http"
)

func HandleFuncs() {

	patternHandlers := map[string]http.Handler{
		// start at the account
		"/":            http.RedirectHandler("/account", http.StatusPermanentRedirect),
		"/account":     memCache(nod.RequestLog(http.HandlerFunc(GetAccount))),
		"/store":       memCache(nod.RequestLog(http.HandlerFunc(GetStore))),
		"/product":     nod.RequestLog(http.HandlerFunc(GetProduct)),
		"/search":      nod.RequestLog(http.HandlerFunc(GetSearch)),
		"/images":      nod.RequestLog(http.HandlerFunc(GetImages)),
		"/videos":      nod.RequestLog(http.HandlerFunc(GetVideos)),
		"/files":       basicHttpAuth(nod.RequestLog(http.HandlerFunc(GetFiles))),
		"/local-file/": basicHttpAuth(nod.RequestLog(http.HandlerFunc(GetLocalFile))),
	}

	for p, h := range patternHandlers {
		http.HandleFunc(p, h.ServeHTTP)
	}
}

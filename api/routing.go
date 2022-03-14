package api

import (
	"github.com/boggydigital/nod"
	"net/http"
)

func HandleFuncs() {

	patternHandlers := map[string]func(http.ResponseWriter, *http.Request){
		// start at the account
		"/":            http.RedirectHandler("/account", http.StatusPermanentRedirect).ServeHTTP,
		"/account":     nod.RequestLog(GetAccount),
		"/store":       nod.RequestLog(GetStore),
		"/product":     nod.RequestLog(GetProduct),
		"/search":      nod.RequestLog(GetSearch),
		"/images":      nod.RequestLog(GetImages),
		"/videos":      nod.RequestLog(GetVideos),
		"/files":       nod.RequestLog(GetFiles),
		"/local-file/": nod.RequestLog(GetLocalFile),
	}

	for p, h := range patternHandlers {
		http.HandleFunc(p, h)
	}
}

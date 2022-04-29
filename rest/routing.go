package rest

import (
	"github.com/boggydigital/nod"
	"net/http"
)

func HandleFuncs() {

	patternHandlers := map[string]http.Handler{
		// current endpoints
		"/updates":     Gzip(nod.RequestLog(http.HandlerFunc(GetUpdates))),
		"/product":     Gzip(nod.RequestLog(http.HandlerFunc(GetProduct))),
		"/search":      Gzip(nod.RequestLog(http.HandlerFunc(GetSearch))),
		"/images":      nod.RequestLog(http.HandlerFunc(GetImages)),
		"/videos":      nod.RequestLog(http.HandlerFunc(GetVideos)),
		"/items/":      nod.RequestLog(http.HandlerFunc(GetItems)),
		"/files":       basicHttpAuth(nod.RequestLog(http.HandlerFunc(GetFiles))),
		"/local-file/": basicHttpAuth(nod.RequestLog(http.HandlerFunc(GetLocalFile))),
		"/favicon.ico": http.HandlerFunc(http.NotFound),

		// updates redirects
		"/updates/recent":    http.RedirectHandler("/updates?since=4", http.StatusPermanentRedirect),
		"/updates/today":     http.RedirectHandler("/updates?since=24", http.StatusPermanentRedirect),
		"/updates/this_week": http.RedirectHandler("/updates?since=120", http.StatusPermanentRedirect),

		// products redirects
		"/products/downloads": http.RedirectHandler(
			"/search?scope=downloads&types=account-products&sort=gog-order-date&desc=true",
			http.StatusPermanentRedirect),
		"/products/wishlist": http.RedirectHandler(
			"/search?scope=wishlist&wishlisted=true&sort=rating&desc=true",
			http.StatusPermanentRedirect),
		"/products/store": http.RedirectHandler(
			"/search?scope=store&types=store-products&sort=gog-release-date&desc=true",
			http.StatusPermanentRedirect),
		"/products": http.RedirectHandler("/search", http.StatusPermanentRedirect),

		// start at the account
		"/": http.RedirectHandler("/updates", http.StatusPermanentRedirect),
	}

	for p, h := range patternHandlers {
		http.HandleFunc(p, h.ServeHTTP)
	}
}

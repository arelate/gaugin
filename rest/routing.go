package rest

import (
	"github.com/boggydigital/middleware"
	"github.com/boggydigital/nod"
	"net/http"
)

var predefinedSearchPaths = map[string]string{
	"owned":    "/search?types=account-products&sort=gog-order-date&desc=true",
	"wishlist": "/search?wishlisted=true&sort=gog-release-date&desc=true",
	"sale":     "/search?types=store-products&owned=false&is-discounted=true&sort=discount-percentage&desc=true",
	"all":      "/search?types=store-products&sort=gog-release-date&desc=true",
}

func HandleFuncs() {

	patternHandlers := map[string]http.Handler{
		// current endpoints
		"/updates":        middleware.Gzip(middleware.GetMethodOnly(nod.RequestLog(http.HandlerFunc(GetUpdates)))),
		"/product":        middleware.Gzip(middleware.GetMethodOnly(nod.RequestLog(http.HandlerFunc(GetProduct)))),
		"/search":         middleware.Gzip(middleware.GetMethodOnly(nod.RequestLog(http.HandlerFunc(GetSearch)))),
		"/description":    middleware.Gzip(middleware.GetMethodOnly(nod.RequestLog(http.HandlerFunc(GetDescription)))),
		"/downloads":      middleware.Gzip(middleware.GetMethodOnly(nod.RequestLog(http.HandlerFunc(GetDownloads)))),
		"/changelog":      middleware.Gzip(middleware.GetMethodOnly(nod.RequestLog(http.HandlerFunc(GetChangelog)))),
		"/screenshots":    middleware.Gzip(middleware.GetMethodOnly(nod.RequestLog(http.HandlerFunc(GetScreenshots)))),
		"/videos":         middleware.Gzip(middleware.GetMethodOnly(nod.RequestLog(http.HandlerFunc(GetVideos)))),
		"/steam-app-news": middleware.Gzip(middleware.GetMethodOnly(nod.RequestLog(http.HandlerFunc(GetSteamAppNews)))),
		"/image":          middleware.GetMethodOnly(nod.RequestLog(http.HandlerFunc(GetImage))),
		"/video":          middleware.GetMethodOnly(nod.RequestLog(http.HandlerFunc(GetVideo))),
		"/thumbnails":     middleware.GetMethodOnly(nod.RequestLog(http.HandlerFunc(GetThumbnails))),
		"/items/":         middleware.GetMethodOnly(nod.RequestLog(http.HandlerFunc(GetItems))),
		"/files":          middleware.BasicHttpAuth(middleware.GetMethodOnly(nod.RequestLog(http.HandlerFunc(GetFiles)))),
		"/local-file/":    middleware.BasicHttpAuth(middleware.GetMethodOnly(nod.RequestLog(http.HandlerFunc(GetLocalFile)))),
		"/favicon.ico":    http.HandlerFunc(http.NotFound),

		// updates redirects
		"/updates/recent":    http.RedirectHandler("/updates?since=4", http.StatusPermanentRedirect),
		"/updates/today":     http.RedirectHandler("/updates?since=24", http.StatusPermanentRedirect),
		"/updates/this_week": http.RedirectHandler("/updates?since=120", http.StatusPermanentRedirect),

		// products redirects
		"/products/owned":    http.RedirectHandler(predefinedSearchPaths["owned"], http.StatusPermanentRedirect),
		"/products/wishlist": http.RedirectHandler(predefinedSearchPaths["wishlist"], http.StatusPermanentRedirect),
		"/products/sale":     http.RedirectHandler(predefinedSearchPaths["sale"], http.StatusPermanentRedirect),
		"/products/all":      http.RedirectHandler(predefinedSearchPaths["all"], http.StatusPermanentRedirect),
		"/products":          http.RedirectHandler("/search", http.StatusPermanentRedirect),

		// start at the updates
		"/": http.RedirectHandler("/updates", http.StatusPermanentRedirect),
	}

	for p, h := range patternHandlers {
		http.HandleFunc(p, h.ServeHTTP)
	}
}

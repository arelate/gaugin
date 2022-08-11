package rest

import (
	"github.com/boggydigital/middleware"
	"github.com/boggydigital/nod"
	"net/http"
)

var (
	Auth     = middleware.BasicHttpAuth
	GetOnly  = middleware.GetMethodOnly
	PostOnly = middleware.PostMethodOnly
	Gzip     = middleware.Gzip
	Log      = nod.RequestLog
	Redirect = http.RedirectHandler
)

var searchRoutes = map[string]string{
	"owned":    "/search?types=account-products&sort=gog-order-date&desc=true",
	"wishlist": "/search?wishlisted=true&sort=gog-release-date&desc=true",
	"sale":     "/search?types=store-products&owned=false&is-discounted=true&sort=discount-percentage&desc=true",
	"all":      "/search?types=store-products&sort=gog-release-date&desc=true",
}

func HandleFuncs() {

	patternHandlers := map[string]http.Handler{
		// unauthenticated endpoints
		"/updates":     Gzip(GetOnly(Log(http.HandlerFunc(GetUpdates)))),
		"/product":     Gzip(Log(http.HandlerFunc(GetProduct))), // can be GET or POST (/tag/apply redirect)
		"/search":      Gzip(GetOnly(Log(http.HandlerFunc(GetSearch)))),
		"/description": Gzip(GetOnly(Log(http.HandlerFunc(GetDescription)))),
		"/downloads":   Gzip(GetOnly(Log(http.HandlerFunc(GetDownloads)))),
		"/changelog":   Gzip(GetOnly(Log(http.HandlerFunc(GetChangelog)))),
		"/screenshots": Gzip(GetOnly(Log(http.HandlerFunc(GetScreenshots)))),
		"/videos":      Gzip(GetOnly(Log(http.HandlerFunc(GetVideos)))),
		"/steam-news":  Gzip(GetOnly(Log(http.HandlerFunc(GetSteamNews)))),
		"/image":       GetOnly(Log(http.HandlerFunc(GetImage))),
		"/video":       GetOnly(Log(http.HandlerFunc(GetVideo))),
		"/thumbnails":  GetOnly(Log(http.HandlerFunc(GetThumbnails))),
		"/items/":      GetOnly(Log(http.HandlerFunc(GetItems))),

		// authenticated endpoints
		"/files":           Auth(GetOnly(Log(http.HandlerFunc(GetFiles)))),
		"/local-file/":     Auth(GetOnly(Log(http.HandlerFunc(GetLocalFile)))),
		"/wishlist/add":    Auth(GetOnly(Log(http.HandlerFunc(GetWishlistAdd)))),
		"/wishlist/remove": Auth(GetOnly(Log(http.HandlerFunc(GetWishlistRemove)))),
		"/tags/edit":       Auth(GetOnly(Log(http.HandlerFunc(GetTagsEdit)))),
		"/tags/apply":      Auth(PostOnly(Log(http.HandlerFunc(PostTagsApply)))),

		// products redirects
		"/products": Redirect("/search", http.StatusPermanentRedirect),

		// start at the updates
		"/": Redirect("/updates", http.StatusPermanentRedirect),
	}

	for route, path := range searchRoutes {
		patternHandlers["/products/"+route] = Redirect(path, http.StatusPermanentRedirect)
	}

	for p, h := range patternHandlers {
		http.HandleFunc(p, h.ServeHTTP)
	}
}

func unconstrainedPath(p string) string {
	return p + "&unconstrained"
}

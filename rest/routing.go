package rest

import (
	"github.com/boggydigital/middleware"
	"github.com/boggydigital/nod"
	"net/http"
)

var (
	BHA = middleware.BasicHttpAuth
	GMO = middleware.GetMethodOnly
	PMO = middleware.PostMethodOnly
	GZ  = middleware.Gzip
	LOG = nod.RequestLog
)

var predefinedSearchPaths = map[string]string{
	"owned":    "/search?types=account-products&sort=gog-order-date&desc=true",
	"wishlist": "/search?wishlisted=true&sort=gog-release-date&desc=true",
	"sale":     "/search?types=store-products&owned=false&is-discounted=true&sort=discount-percentage&desc=true",
	"all":      "/search?types=store-products&sort=gog-release-date&desc=true",
}

func HandleFuncs() {

	patternHandlers := map[string]http.Handler{
		// unauthenticated endpoints
		"/updates":        GZ(GMO(LOG(http.HandlerFunc(GetUpdates)))),
		"/product":        GZ(LOG(http.HandlerFunc(GetProduct))), // can be GET or POST (/tag/apply redirect)
		"/search":         GZ(GMO(LOG(http.HandlerFunc(GetSearch)))),
		"/description":    GZ(GMO(LOG(http.HandlerFunc(GetDescription)))),
		"/downloads":      GZ(GMO(LOG(http.HandlerFunc(GetDownloads)))),
		"/changelog":      GZ(GMO(LOG(http.HandlerFunc(GetChangelog)))),
		"/screenshots":    GZ(GMO(LOG(http.HandlerFunc(GetScreenshots)))),
		"/videos":         GZ(GMO(LOG(http.HandlerFunc(GetVideos)))),
		"/steam-app-news": GZ(GMO(LOG(http.HandlerFunc(GetSteamAppNews)))),
		"/image":          GMO(LOG(http.HandlerFunc(GetImage))),
		"/video":          GMO(LOG(http.HandlerFunc(GetVideo))),
		"/thumbnails":     GMO(LOG(http.HandlerFunc(GetThumbnails))),
		"/items/":         GMO(LOG(http.HandlerFunc(GetItems))),
		"/favicon.ico":    http.HandlerFunc(http.NotFound),

		// authenticated endpoints
		"/files":           BHA(GMO(LOG(http.HandlerFunc(GetFiles)))),
		"/local-file/":     BHA(GMO(LOG(http.HandlerFunc(GetLocalFile)))),
		"/wishlist/add":    BHA(GMO(LOG(http.HandlerFunc(GetWishlistAdd)))),
		"/wishlist/remove": BHA(GMO(LOG(http.HandlerFunc(GetWishlistRemove)))),
		"/tags/edit":       BHA(GMO(LOG(http.HandlerFunc(GetTagsEdit)))),
		"/tags/apply":      BHA(PMO(LOG(http.HandlerFunc(PostTagsApply)))),

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

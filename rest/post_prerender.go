package rest

import (
	"github.com/boggydigital/nod"
	"github.com/boggydigital/stencil/stencil_rest"
	"net/http"
)

func PostPrerender(w http.ResponseWriter, _ *http.Request) {

	// POST /prerender

	// the following pages will be pre-rendered:
	// - default path (/updates)
	// - every top-level search route (/search, owned, wishlist, sale, all)
	// - every product updated at the last sync

	if err := updatePrerender(); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}
}

func updatePrerender() error {

	ps := make([]string, 0)
	ps = append(ps, "/updates")

	for _, p := range searchRoutes() {
		ps = append(ps, p)
	}

	return stencil_rest.Prerender(ps, port)
}

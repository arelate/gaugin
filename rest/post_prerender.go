package rest

import (
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/stencil/stencil_rest"
	"net/http"
)

func PostPrerender(w http.ResponseWriter, r *http.Request) {

	// POST /prerender

	// the following pages will be pre-rendered:
	// - default path (/updates)
	// - every top-level search route (/search, owned, wishlist, sale, all)
	// - every product updated at the last sync

	paths := []string{
		"/updates",
	}

	for _, p := range searchRoutes() {
		paths = append(paths, p)
	}

	updRdx, _, err := getRedux(
		http.DefaultClient,
		"",
		true,
		vangogh_local_data.LastSyncUpdatesProperty)

	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	keys := make(map[string]interface{})
	for _, rdx := range updRdx {
		for _, id := range rdx[vangogh_local_data.LastSyncUpdatesProperty] {
			keys[id] = nil
		}
	}

	for id := range keys {
		paths = append(paths, "/product?id="+id)
	}

	stencil_rest.Prerender(paths, port, w)
}

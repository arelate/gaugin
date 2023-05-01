package rest

import (
	"github.com/arelate/gaugin/paths"
	"github.com/arelate/vangogh_local_data"
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

	if err := setPrerender(); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}
}

func setPrerender() error {
	p := make([]string, 0)
	p = appendListsPaths(p)

	p, err := appendUpdatedItemsPaths(p)
	if err != nil {
		return err
	}

	if err := stencil_rest.Prerender(p, true, port); err != nil {
		return err
	}

	return nil
}

func updatePrerender(ids ...string) error {
	p := make([]string, 0)
	p = appendListsPaths(p)

	for _, id := range ids {
		p = append(p, paths.ProductId(id))
	}

	if err := stencil_rest.Prerender(p, false, port); err != nil {
		return err
	}

	return nil
}

func appendListsPaths(paths []string) []string {
	paths = append(paths, "/updates")

	for _, p := range searchRoutes() {
		paths = append(paths, p)
	}

	return paths
}

func appendUpdatedItemsPaths(p []string) ([]string, error) {

	updRdx, _, err := getRedux(
		http.DefaultClient,
		"",
		true,
		vangogh_local_data.LastSyncUpdatesProperty)

	if err != nil {
		return nil, err
	}

	keys := make(map[string]interface{})
	for _, rdx := range updRdx {
		for _, id := range rdx[vangogh_local_data.LastSyncUpdatesProperty] {
			keys[id] = nil
		}
	}

	for id := range keys {
		p = append(p, paths.ProductId(id))
	}

	return p, nil
}

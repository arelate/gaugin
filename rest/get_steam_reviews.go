package rest

import (
	"github.com/arelate/gaugin/rest/compton_pages"
	"github.com/boggydigital/nod"
	"net/http"
)

func GetSteamReviews(w http.ResponseWriter, r *http.Request) {

	// GET /steam-reviews?id

	id := r.URL.Query().Get("id")

	sar, err := getSteamReviews(http.DefaultClient, id)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	p := compton_pages.SteamReviews(id, sar)

	if err := p.WriteContent(w); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
	}
}

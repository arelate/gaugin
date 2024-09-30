package rest

import (
	"github.com/arelate/gaugin/rest/compton_pages"
	"github.com/boggydigital/nod"
	"net/http"
)

func GetSteamNews(w http.ResponseWriter, r *http.Request) {

	// GET /steam-news?id

	id := r.URL.Query().Get("id")
	all := r.URL.Query().Has("all")

	san, err := getSteamNews(http.DefaultClient, id)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	p := compton_pages.SteamNews(id, san, all)

	if err := p.WriteContent(w); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
	}
}

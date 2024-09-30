package rest

import (
	"github.com/arelate/gaugin/rest/compton_pages"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/kevlar"
	"github.com/boggydigital/nod"
	"net/http"
)

func GetSteamDeck(w http.ResponseWriter, r *http.Request) {

	// GET /steam-deck?id

	id := r.URL.Query().Get("id")

	dacr, err := getSteamDeckReport(http.DefaultClient, id)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	idRedux, err := getRedux(http.DefaultClient, id, false,
		vangogh_local_data.TitleProperty,
		vangogh_local_data.SteamDeckAppCompatibilityCategoryProperty,
	)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	p := compton_pages.SteamDeck(id, dacr, kevlar.ReduxProxy(idRedux))

	if err := p.WriteContent(w); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
	}
}

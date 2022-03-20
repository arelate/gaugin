package rest

import (
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"golang.org/x/exp/maps"
	"net/http"
	"strings"
)

func GetNews(w http.ResponseWriter, r *http.Request) {

	// GET /new?since-hours-ago

	dc := http.DefaultClient

	updates, err := getUpdates(dc, gog_integration.Game, 24*7)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), 405)
		return
	}

	keys := make(map[string]bool)
	for _, ids := range updates {
		for _, id := range ids {
			keys[id] = true
		}
	}

	rdx, err := getRedux(dc,
		strings.Join(maps.Keys(keys), ","),
		vangogh_local_data.TitleProperty,
		vangogh_local_data.DevelopersProperty,
		vangogh_local_data.PublisherProperty)

	if err != nil {
		http.Error(w, "error getting all_redux", http.StatusInternalServerError)
		return
	}

	uvm := updatesViewModelFromRedux(updates, rdx)

	if err := tmpl.ExecuteTemplate(w, "news", uvm); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}

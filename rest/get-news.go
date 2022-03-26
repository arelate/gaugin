package rest

import (
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"golang.org/x/exp/maps"
	"net/http"
	"strconv"
	"strings"
)

const defaultSince = 24

func GetUpdates(w http.ResponseWriter, r *http.Request) {

	// GET /new?since-hours-ago

	dc := http.DefaultClient

	sinceStr := vangogh_local_data.ValueFromUrl(r.URL, "since")
	since, err := strconv.Atoi(sinceStr)
	if err != nil {
		since = 0
	}
	if since == 0 {
		since = defaultSince
	}

	updates, err := getUpdates(dc, gog_integration.Game, since)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusMethodNotAllowed)
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
		vangogh_local_data.WishlistedProperty,
		vangogh_local_data.DevelopersProperty,
		vangogh_local_data.PublisherProperty,
		vangogh_local_data.OperatingSystemsProperty,
		vangogh_local_data.TagIdProperty,
		vangogh_local_data.ProductTypeProperty)

	if err != nil {
		http.Error(w, "error getting all_redux", http.StatusInternalServerError)
		return
	}

	uvm := updatesViewModelFromRedux(updates, since, rdx)

	if err := tmpl.ExecuteTemplate(w, "updates", uvm); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}

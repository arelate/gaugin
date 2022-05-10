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

var listReduxProperties = []string{
	vangogh_local_data.TitleProperty,
	vangogh_local_data.WishlistedProperty,
	vangogh_local_data.OwnedProperty,
	vangogh_local_data.PreOrderProperty,
	vangogh_local_data.ComingSoonProperty,
	vangogh_local_data.TBAProperty,
	vangogh_local_data.InDevelopmentProperty,
	vangogh_local_data.IsUsingDOSBoxProperty,
	vangogh_local_data.IsUsingScummVMProperty,
	vangogh_local_data.DevelopersProperty,
	vangogh_local_data.PublisherProperty,
	vangogh_local_data.OperatingSystemsProperty,
	vangogh_local_data.TagIdProperty,
	vangogh_local_data.ProductTypeProperty}

func GetUpdates(w http.ResponseWriter, r *http.Request) {

	// GET /updates?since=hours-ago

	sinceStr := vangogh_local_data.ValueFromUrl(r.URL, "since")
	since, err := strconv.Atoi(sinceStr)
	if err != nil {
		since = 0
	}
	if since <= 0 {
		since = defaultSince
	}

	updates, err := getUpdates(http.DefaultClient, gog_integration.Game, since)
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

	rdx, err := getRedux(
		http.DefaultClient,
		strings.Join(maps.Keys(keys), ","),
		listReduxProperties...)

	if err != nil {
		http.Error(w, nod.ErrorStr("error getting all_redux"), http.StatusInternalServerError)
		return
	}

	uvm := updatesViewModelFromRedux(updates, since, rdx)

	if err := tmpl.ExecuteTemplate(w, "updates", uvm); err != nil {
		http.Error(w, nod.ErrorStr("template exec error"), http.StatusInternalServerError)
		return
	}
}

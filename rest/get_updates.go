package rest

import (
	"net/http"
	"strings"
	"time"

	"github.com/arelate/gaugin/gaugin_middleware"
	"github.com/arelate/gaugin/view_models"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"golang.org/x/exp/maps"
)

func GetUpdates(w http.ResponseWriter, r *http.Request) {

	// GET /updates

	timing := make(map[string]int64)

	start := time.Now()

	updRdx, err := getRedux(
		http.DefaultClient,
		"",
		true,
		vangogh_local_data.LastSyncUpdatesProperty)

	timing["updRdx"] = time.Since(start).Milliseconds()

	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	updates := make(map[string][]string)
	for section, rdx := range updRdx {
		ids := rdx[vangogh_local_data.LastSyncUpdatesProperty]
		if len(ids) > 0 {
			updates[section] = ids
		}
	}

	keys := make(map[string]bool)
	for _, ids := range updates {
		for _, id := range ids {
			keys[id] = true
		}
	}

	start = time.Now()

	dataRdx, err := getRedux(
		http.DefaultClient,
		strings.Join(maps.Keys(keys), ","),
		false,
		view_models.ListProperties...)

	timing["dataRdx"] = time.Since(start).Milliseconds()

	if err != nil {
		http.Error(w, nod.ErrorStr("error getting redux"), http.StatusInternalServerError)
		return
	}

	start = time.Now()

	syncRdx, err := getRedux(
		http.DefaultClient,
		vangogh_local_data.SyncCompleteKey,
		false,
		vangogh_local_data.SyncEventsProperty)

	timing["syncRdx"] = time.Since(start).Milliseconds()

	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	uvm := view_models.NewUpdates(updates, dataRdx, syncRdx[vangogh_local_data.SyncCompleteKey])

	gaugin_middleware.DefaultHeaders(timing, w)

	if err := tmpl.ExecuteTemplate(w, "updates-page", uvm); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}
}

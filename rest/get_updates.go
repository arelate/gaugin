package rest

import (
	"net/http"
	"sort"
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

	st := gaugin_middleware.NewServerTimings()

	start := time.Now()
	updRdx, cached, err := getRedux(
		http.DefaultClient,
		"",
		true,
		vangogh_local_data.LastSyncUpdatesProperty)

	if cached {
		st.SetFlag("updRdx-cached")
	}
	st.Set("updRdx", time.Since(start).Milliseconds())

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

	ids := maps.Keys(keys)
	sort.Strings(ids)

	start = time.Now()
	dataRdx, cached, err := getRedux(
		http.DefaultClient,
		strings.Join(ids, ","),
		false,
		view_models.ListProperties...)

	if cached {
		st.SetFlag("dataRdx-cached")
	}
	st.Set("dataRdx", time.Since(start).Milliseconds())

	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	start = time.Now()
	syncRdx, cached, err := getRedux(
		http.DefaultClient,
		vangogh_local_data.SyncCompleteKey,
		false,
		vangogh_local_data.SyncEventsProperty)

	if cached {
		st.SetFlag("syncRdx-cached")
	}
	st.Set("syncRdx", time.Since(start).Milliseconds())

	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	uvm := view_models.NewUpdates(updates, dataRdx, syncRdx[vangogh_local_data.SyncCompleteKey])

	gaugin_middleware.DefaultHeaders(st, w)

	if err := tmpl.ExecuteTemplate(w, "updates-page", uvm); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}
}

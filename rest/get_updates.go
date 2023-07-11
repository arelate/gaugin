package rest

import (
	"github.com/arelate/gaugin/stencil_app"
	"golang.org/x/exp/maps"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/arelate/gaugin/gaugin_middleware"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
)

const (
	updatedProductsLimit = 24
)

func GetUpdates(w http.ResponseWriter, r *http.Request) {

	// GET /updates

	showAll := r.URL.Query().Get("show-all") == "true"

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
	updateTotals := make(map[string]int)

	for section, rdx := range updRdx {
		ids := rdx[vangogh_local_data.LastSyncUpdatesProperty]
		updateTotals[section] = len(ids)
		for _, id := range ids {
			if !showAll && len(updates[section]) >= updatedProductsLimit {
				continue
			}
			updates[section] = append(updates[section], id)
		}
		//if len(ids) > 0 {
		//	updates[section] = ids
		//}
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
		stencil_app.ProductsProperties...)

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

	updated := "recently"
	syncDra := vangogh_local_data.NewIRAProxy(syncRdx)
	if scs, ok := syncDra.GetFirstVal(vangogh_local_data.SyncEventsProperty, vangogh_local_data.SyncCompleteKey); ok {
		if sci, err := strconv.ParseInt(scs, 10, 64); err == nil {
			updated = time.Unix(sci, 0).Format(time.RFC1123)
		}
	}

	irap := vangogh_local_data.NewIRAProxy(dataRdx)

	gaugin_middleware.DefaultHeaders(st, w)

	sections := maps.Keys(updates)
	sort.Strings(sections)

	var caser = cases.Title(language.Russian)

	sectionTitles := make(map[string]string)
	for t, _ := range updates {
		sectionTitles[t] = caser.String(t)
	}

	if err := app.RenderGroup(
		stencil_app.NavUpdates,
		sections,
		updates,
		sectionTitles,
		updateTotals,
		updated,
		r.URL,
		irap,
		w); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}
}

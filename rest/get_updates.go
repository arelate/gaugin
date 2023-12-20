package rest

import (
	"github.com/arelate/gaugin/stencil_app"
	"github.com/boggydigital/kvas"
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

	paginate := false

	for section, rdx := range updRdx {
		ids := rdx[vangogh_local_data.LastSyncUpdatesProperty]
		updateTotals[section] = len(ids)
		// limit number of items only if there are at least x2 the limit
		// e.g. if the limit is 24, only start limiting if there are 49 or more items
		paginate = len(ids) > updatedProductsLimit*2
		for _, id := range ids {
			if paginate && !showAll && len(updates[section]) >= updatedProductsLimit {
				continue
			}
			updates[section] = append(updates[section], id)
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
	syncDra := kvas.ReduxProxy(syncRdx)
	if scs, ok := syncDra.GetFirstVal(vangogh_local_data.SyncEventsProperty, vangogh_local_data.SyncCompleteKey); ok {
		if sci, err := strconv.ParseInt(scs, 10, 64); err == nil {
			updated = time.Unix(sci, 0).Format(time.RFC1123)
		}
	}

	gaugin_middleware.DefaultHeaders(st, w)

	// section order will be based on full title ("new in ...", "updates in ...")
	// so not changed after concise update titles change
	sections := maps.Keys(updates)
	sort.Strings(sections)

	var caser = cases.Title(language.English)

	sectionTitles := make(map[string]string)
	for t, _ := range updates {
		st := t
		switch t {
		case "new in store":
			st = "store"
		case "new in account":
			st = "account"
		case "new in wishlist":
			st = "wishlist"
		case "released today":
			st = "today"
		case "updates in account":
			st = "updates"
		case "updates in news":
			st = "news"
		}

		sectionTitles[t] = caser.String(st)
	}

	tagNamesRedux, _, err := getRedux(http.DefaultClient, "", true, vangogh_local_data.TagNameProperty)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	rdx := kvas.ReduxProxy(MergeIdPropertyValues(dataRdx, tagNamesRedux))

	if err := app.RenderGroup(
		stencil_app.NavUpdates,
		sections,
		updates,
		sectionTitles,
		updateTotals,
		updated,
		r.URL,
		rdx,
		w); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}
}

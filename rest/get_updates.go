package rest

import (
	"github.com/arelate/gaugin/rest/compton_data"
	"github.com/arelate/gaugin/rest/compton_pages"
	"github.com/boggydigital/kevlar"
	"golang.org/x/exp/maps"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
)

const (
	updatedProductsLimit = 24 // divisible by 2,3,4,6
)

var sectionTitles = map[string]string{
	"new in store":       "Store additions",
	"new in account":     "Purchased recently",
	"new in wishlist":    "Wishlist additions",
	"released today":     "Today's releases",
	"updates in account": "Updated installers",
	"updates in news":    "Steam news",
}

func GetUpdates(w http.ResponseWriter, r *http.Request) {

	// GET /updates

	showAll := r.URL.Query().Get("show-all") == "true"

	updRdx, err := getRedux(
		http.DefaultClient,
		"",
		true,
		vangogh_local_data.LastSyncUpdatesProperty)

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

	dataRdx, err := getRedux(
		http.DefaultClient,
		strings.Join(ids, ","),
		false,
		compton_data.ProductsProperties...)

	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	syncRdx, err := getRedux(
		http.DefaultClient,
		vangogh_local_data.SyncCompleteKey,
		false,
		vangogh_local_data.SyncEventsProperty)

	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	updated := "recently"
	syncDra := kevlar.ReduxProxy(syncRdx)
	if scs, ok := syncDra.GetLastVal(vangogh_local_data.SyncEventsProperty, vangogh_local_data.SyncCompleteKey); ok {
		if sci, err := strconv.ParseInt(scs, 10, 64); err == nil {
			updated = time.Unix(sci, 0).Format(time.RFC1123)
		}
	}

	// section order will be based on full title ("new in ...", "updates in ...")
	// so the order won't be changed after expanding titles
	sections := maps.Keys(updates)
	sort.Strings(sections)

	tagNamesRedux, err := getRedux(http.DefaultClient, "", true, vangogh_local_data.TagNameProperty)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	rdx := kevlar.ReduxProxy(MergeIdPropertyValues(dataRdx, tagNamesRedux))

	updatesPage := compton_pages.Updates(sections, updates, sectionTitles, updateTotals, updated, rdx)
	if err := updatesPage.WriteContent(w); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
	}
}

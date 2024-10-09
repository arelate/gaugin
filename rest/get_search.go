package rest

import (
	"github.com/arelate/gaugin/rest/compton_data"
	"github.com/arelate/gaugin/rest/compton_pages"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/kevlar"
	"github.com/boggydigital/middleware"
	"github.com/boggydigital/nod"
	"net/http"
	"strconv"
	"strings"
)

func GetSearch(w http.ResponseWriter, r *http.Request) {

	// GET /search?(search_params)&from

	q := r.URL.Query()

	from, to := 0, 0
	if q.Has("from") {
		from64, err := strconv.ParseInt(q.Get("from"), 10, 32)
		if err != nil {
			http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
			return
		}
		from = int(from64)
	}

	query := make(map[string][]string)

	shortQuery := false
	queryProperties := compton_data.SearchProperties
	for _, p := range queryProperties {
		if v := q.Get(p); v != "" {
			query[p] = strings.Split(v, ",")
		} else {
			if q.Has(p) {
				q.Del(p)
				shortQuery = true
			}
		}
	}

	//if we removed some properties with no values - redirect to the shortest URL
	if shortQuery {
		r.URL.RawQuery = q.Encode()
		http.Redirect(w, r, r.URL.String(), http.StatusPermanentRedirect)
		return
	}

	var ids, slice []string

	dc := http.DefaultClient

	idRedux := NewIdPropertyValues(compton_data.ProductsProperties)

	if len(query) > 0 {

		var err error
		ids, err = getSearch(dc, q)
		if err != nil {
			http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
			return
		}

		if from > len(ids)-1 {
			from = 0
		}

		to = from + SearchResultsLimit
		if to > len(ids) {
			to = len(ids)
		} else if to+SearchResultsLimit > len(ids) {
			to = len(ids)
		}

		slice = ids[from:to]

		su := searchUrl(q)

		lm := urlLastModified[su.String()]
		w.Header().Set(middleware.LastModifiedHeader, lm)

		ims := r.Header.Get(middleware.IfModifiedSinceHeader)
		if middleware.IsNotModified(ims, lm) {
			w.WriteHeader(http.StatusNotModified)
			return
		}

		irx, err := getRedux(dc, strings.Join(slice, ","), false, compton_data.ProductsProperties...)
		if err != nil {
			http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
			return
		}

		idRedux = irx
	}

	// adding tag names for related games

	tagNamesRedux, err := getRedux(http.DefaultClient, "", true, vangogh_local_data.TagNameProperty)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	idRedux = MergeIdPropertyValues(idRedux, tagNamesRedux)

	rdx := kevlar.ReduxProxy(idRedux)

	searchPage := compton_pages.Search(query, ids, from, to, rdx)
	if err := searchPage.WriteContent(w); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
	}
}

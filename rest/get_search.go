package rest

import (
	"github.com/arelate/gaugin/stencil_app"
	"github.com/boggydigital/kevlar"
	"github.com/boggydigital/middleware"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/arelate/gaugin/gaugin_middleware"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
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
	queryProperties := stencil_app.SearchProperties
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

	var start time.Time
	st := gaugin_middleware.NewServerTimings()

	idRedux := NewIdPropertyValues(stencil_app.ProductsProperties)

	if len(query) > 0 {
		start = time.Now()

		var cached bool
		var err error
		ids, cached, err = getSearch(dc, q)
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

		if cached {
			st.SetFlag("getSearch-cached")
		}
		st.Set("getSearch", time.Since(start).Milliseconds())

		su := searchUrl(q)

		lm := urlLastModified[su.String()]
		w.Header().Set(middleware.LastModifiedHeader, lm)

		ims := r.Header.Get(middleware.IfModifiedSinceHeader)
		if middleware.IsNotModified(ims, lm) {
			w.WriteHeader(http.StatusNotModified)
			return
		}

		start = time.Now()
		irx, cached, err := getRedux(dc, strings.Join(slice, ","), false, stencil_app.ProductsProperties...)
		if err != nil {
			http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
			return
		}

		if cached {
			st.SetFlag("getRedux-cached")
		}
		st.Set("getRedux", time.Since(start).Milliseconds())

		// adding tag names for related games

		tagNamesRedux, _, err := getRedux(http.DefaultClient, "", true, vangogh_local_data.TagNameProperty)
		if err != nil {
			http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
			return
		}

		if len(irx) > 0 {
			idRedux = MergeIdPropertyValues(irx, tagNamesRedux)
		}
	}

	rdx := kevlar.ReduxProxy(idRedux)

	gaugin_middleware.DefaultHeaders(st, w)

	if err := app.RenderSearch(stencil_app.NavSearch, query, slice, from, to, len(ids), r.URL, rdx, w); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}
}

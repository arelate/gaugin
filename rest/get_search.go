package rest

import (
	"github.com/arelate/gaugin/stencil_app"
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

	irap := vangogh_local_data.NewEmptyIRAProxy(stencil_app.ProductsProperties)

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
		}

		slice = ids[from:to]

		if cached {
			st.SetFlag("getSearch-cached")
		}
		st.Set("getSearch", time.Since(start).Milliseconds())

		su := searchUrl(q)

		lmu := time.Unix(urlLastModified[su.String()], 0).UTC()
		w.Header().Set("Last-Modified", lmu.Format(time.RFC1123))

		if ims := r.Header.Get("If-Modified-Since"); ims != "" {
			if imst, err := time.Parse(time.RFC1123, ims); err == nil {
				if imst.UTC().Unix() <= lmu.Unix() {
					w.WriteHeader(http.StatusNotModified)
					return
				}
			}
		}

		start = time.Now()
		rdx, cached, err := getRedux(dc, strings.Join(slice, ","), false, stencil_app.ProductsProperties...)
		if err != nil {
			http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
			return
		}

		irap = vangogh_local_data.NewIRAProxy(rdx)

		if cached {
			st.SetFlag("getRedux-cached")
		}
		st.Set("getRedux", time.Since(start).Milliseconds())
	}

	gaugin_middleware.DefaultHeaders(st, w)

	if err := app.RenderSearch(stencil_app.NavSearch, query, slice, from, to, len(ids), r.URL, irap, w); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}
}

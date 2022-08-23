package rest

import (
	"net/http"
	"strings"
	"time"

	"github.com/arelate/gaugin/gaugin_middleware"
	"github.com/arelate/gaugin/view_models"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
)

func GetSearch(w http.ResponseWriter, r *http.Request) {

	// GET /search?(search_params)

	scope := ""
	for route, path := range searchRoutes() {
		if r.URL.RawQuery != "" &&
			(strings.HasSuffix(path, r.URL.RawQuery) ||
				strings.HasSuffix(unconstrainedPath(path), r.URL.RawQuery)) {
			scope = route
			break
		}
		if r.URL.RawQuery == "" &&
			r.URL.Path == path {
			scope = route
			break
		}
	}

	q := r.URL.Query()

	constrained := !vangogh_local_data.FlagFromUrl(r.URL, "unconstrained")
	path := ""
	if constrained {
		path = r.URL.RawPath + "?" + r.URL.RawQuery
	}

	spvm := view_models.NewSearchProducts(
		scope,
		constrained,
		path,
	)

	shortQuery := false
	queryProperties := view_models.SearchProperties
	for _, p := range queryProperties {
		if v := q.Get(p); v != "" {
			spvm.Query[p] = v
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

	dc := http.DefaultClient

	var start time.Time
	st := gaugin_middleware.NewServerTimings()

	if len(spvm.Query) > 0 {
		start = time.Now()
		keys, cached, err := getSearch(dc, q)
		if err != nil {
			http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
			return
		}

		if cached {
			st.SetFlag("getSearch-cached")
		}
		st.Set("getSearch", time.Since(start).Milliseconds())

		spvm.Total = len(keys)
		if spvm.Total > spvm.Limit && spvm.Constrained {
			keys = keys[:view_models.SearchProductsLimit]
		}

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
		rdx, cached, err := getRedux(dc, strings.Join(keys, ","), false, view_models.ListProperties...)

		if err != nil {
			http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
			return
		}

		if cached {
			st.SetFlag("getRedux-cached")
		}
		st.Set("getRedux", time.Since(start).Milliseconds())

		lvm := view_models.NewListViewModel(keys, rdx)
		spvm.Products = lvm.Products
	}

	// checking outside search action to account for empty query case
	if spvm.Total <= spvm.Limit {
		spvm.Constrained = false
	}

	start = time.Now()
	digests, cached, err := getDigests(dc, view_models.DigestProperties...)

	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	if cached {
		st.SetFlag("getDigests-cached")
	}
	st.Set("getDigests", time.Since(start).Milliseconds())

	digests[vangogh_local_data.SortProperty] = []string{
		vangogh_local_data.GlobalReleaseDateProperty,
		vangogh_local_data.GOGReleaseDateProperty,
		vangogh_local_data.GOGOrderDateProperty,
		vangogh_local_data.TitleProperty,
		vangogh_local_data.RatingProperty,
		vangogh_local_data.DiscountPercentageProperty}

	digests[vangogh_local_data.DescendingProperty] = []string{
		vangogh_local_data.TrueValue,
		vangogh_local_data.FalseValue}

	spvm.Digests = digests

	gaugin_middleware.DefaultHeaders(st, w)

	if err := tmpl.ExecuteTemplate(w, "search-page", spvm); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}
}

package rest

import (
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/http"
	"strings"
	"time"
)

var gauginSearchProperties = []string{
	vangogh_local_data.TextProperty,
	vangogh_local_data.TitleProperty,
	vangogh_local_data.TagIdProperty,
	vangogh_local_data.OperatingSystemsProperty,
	vangogh_local_data.DevelopersProperty,
	vangogh_local_data.PublisherProperty,
	vangogh_local_data.SeriesProperty,
	vangogh_local_data.GenresProperty,
	vangogh_local_data.PropertiesProperty,
	vangogh_local_data.FeaturesProperty,
	vangogh_local_data.LanguageCodeProperty,
	vangogh_local_data.IncludesGamesProperty,
	vangogh_local_data.IsIncludedByGamesProperty,
	vangogh_local_data.RequiresGamesProperty,
	vangogh_local_data.IsRequiredByGamesProperty,
	vangogh_local_data.ProductTypeProperty,
	vangogh_local_data.WishlistedProperty,
	vangogh_local_data.OwnedProperty,
	vangogh_local_data.IsFreeProperty,
	vangogh_local_data.IsDiscountedProperty,
	vangogh_local_data.PreOrderProperty,
	vangogh_local_data.ComingSoonProperty,
	vangogh_local_data.TBAProperty,
	vangogh_local_data.InDevelopmentProperty,
	vangogh_local_data.IsUsingDOSBoxProperty,
	vangogh_local_data.IsUsingScummVMProperty,
	vangogh_local_data.TypesProperty,
	//
	"sort",
	"desc",
}

var searchPropertyNames = map[string]string{
	vangogh_local_data.TextProperty:              "Any Text",
	vangogh_local_data.TitleProperty:             "Title",
	vangogh_local_data.TagIdProperty:             "User Tags",
	vangogh_local_data.OperatingSystemsProperty:  "OS",
	vangogh_local_data.DevelopersProperty:        "Developers",
	vangogh_local_data.PublisherProperty:         "Publisher",
	vangogh_local_data.SeriesProperty:            "Series",
	vangogh_local_data.GenresProperty:            "Genres",
	vangogh_local_data.PropertiesProperty:        "Store Tags",
	vangogh_local_data.FeaturesProperty:          "Features",
	vangogh_local_data.LanguageCodeProperty:      "Language",
	vangogh_local_data.IncludesGamesProperty:     "Includes",
	vangogh_local_data.IsIncludedByGamesProperty: "Included By",
	vangogh_local_data.RequiresGamesProperty:     "Requires",
	vangogh_local_data.IsRequiredByGamesProperty: "Required By",
	vangogh_local_data.ProductTypeProperty:       "Product Type",
	vangogh_local_data.WishlistedProperty:        "Wishlisted",
	vangogh_local_data.OwnedProperty:             "Owned",
	vangogh_local_data.IsFreeProperty:            "Free",
	vangogh_local_data.IsDiscountedProperty:      "On Sale",
	vangogh_local_data.PreOrderProperty:          "Pre-order",
	vangogh_local_data.ComingSoonProperty:        "Coming Soon",
	vangogh_local_data.TBAProperty:               "TBA",
	vangogh_local_data.InDevelopmentProperty:     "In Development",
	vangogh_local_data.IsUsingDOSBoxProperty:     "Using DOSBox",
	vangogh_local_data.IsUsingScummVMProperty:    "Using ScummVM",
	vangogh_local_data.TypesProperty:             "Data Type",
	//
	"sort": "Sort",
	"desc": "Descending",
}

var gauginDigestibleProperties = []string{
	vangogh_local_data.TagIdProperty,
	vangogh_local_data.OperatingSystemsProperty,
	vangogh_local_data.GenresProperty,
	vangogh_local_data.PropertiesProperty,
	vangogh_local_data.FeaturesProperty,
	vangogh_local_data.LanguageCodeProperty,
	vangogh_local_data.ProductTypeProperty,
	vangogh_local_data.WishlistedProperty,
	vangogh_local_data.OwnedProperty,
	vangogh_local_data.IsFreeProperty,
	vangogh_local_data.IsDiscountedProperty,
	vangogh_local_data.PreOrderProperty,
	vangogh_local_data.ComingSoonProperty,
	vangogh_local_data.TBAProperty,
	vangogh_local_data.InDevelopmentProperty,
	vangogh_local_data.IsUsingDOSBoxProperty,
	vangogh_local_data.IsUsingScummVMProperty,
	vangogh_local_data.TypesProperty,
}

type searchProductsViewModel struct {
	Context          string
	Scope            string
	SearchProperties []string
	Query            map[string]string
	Digests          map[string][]string
	Products         []listProductViewModel
}

func GetSearch(w http.ResponseWriter, r *http.Request) {

	// GET /search?(search_params)

	dc := http.DefaultClient

	scope := ""
	for s, rp := range predefinedSearchPaths {
		if r.URL.RawQuery != "" && strings.HasSuffix(rp, r.URL.RawQuery) {
			scope = s
			break
		}
	}

	q := r.URL.Query()

	spvm := &searchProductsViewModel{
		Scope:            scope,
		Context:          "filter-products",
		SearchProperties: gauginSearchProperties,
		Query:            make(map[string]string, len(q)),
	}

	queryProperties := append(gauginSearchProperties)
	for _, p := range queryProperties {
		if v := q.Get(p); v != "" {
			spvm.Query[p] = v
		}
	}

	if len(spvm.Query) > 0 {
		keys, err := getSearch(dc, q)
		if err != nil {
			http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
			return
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

		rdx, err := getRedux(dc, strings.Join(keys, ","), listReduxProperties...)

		if err != nil {
			http.Error(w, nod.ErrorStr("error getting all_redux"), http.StatusInternalServerError)
			return
		}

		lvm := listViewModelFromRedux(keys, rdx)
		spvm.Products = lvm.Products
	}

	digests, err := getDigests(dc, gauginDigestibleProperties...)

	digests["sort"] = []string{
		vangogh_local_data.GOGReleaseDateProperty,
		vangogh_local_data.GOGOrderDateProperty,
		vangogh_local_data.TitleProperty,
		vangogh_local_data.RatingProperty}

	digests["desc"] = []string{"true", "false"}

	if err != nil {
		http.Error(w, nod.ErrorStr("error getting digests"), http.StatusInternalServerError)
		return
	}
	spvm.Digests = digests

	defaultHeaders(w)

	if err := tmpl.ExecuteTemplate(w, "search", spvm); err != nil {
		http.Error(w, nod.ErrorStr("template error"), http.StatusInternalServerError)
		return
	}
}

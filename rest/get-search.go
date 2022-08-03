package rest

import (
	"github.com/arelate/gaugin/gaugin_middleware"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/http"
	"strings"
	"time"
)

const limit = 100

var gauginSearchProperties = []string{
	vangogh_local_data.TextProperty,
	vangogh_local_data.TitleProperty,
	vangogh_local_data.TagIdProperty,
	vangogh_local_data.LocalTagsProperty,
	vangogh_local_data.OperatingSystemsProperty,
	vangogh_local_data.DevelopersProperty,
	vangogh_local_data.PublisherProperty,
	vangogh_local_data.SeriesProperty,
	vangogh_local_data.GenresProperty,
	vangogh_local_data.PropertiesProperty, //Properties = (GOG.com) Store Tags
	vangogh_local_data.SteamTagsProperty,
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
	vangogh_local_data.SteamReviewScoreDescProperty,
	vangogh_local_data.SortProperty,
	vangogh_local_data.DescendingProperty,
}

var searchPropertyNames = map[string]string{
	vangogh_local_data.TextProperty:                 "Any Text",
	vangogh_local_data.TitleProperty:                "Title",
	vangogh_local_data.TagIdProperty:                "Account Tags",
	vangogh_local_data.LocalTagsProperty:            "Local Tags",
	vangogh_local_data.SteamTagsProperty:            "Steam Tags",
	vangogh_local_data.OperatingSystemsProperty:     "OS",
	vangogh_local_data.DevelopersProperty:           "Developers",
	vangogh_local_data.PublisherProperty:            "Publisher",
	vangogh_local_data.SeriesProperty:               "Series",
	vangogh_local_data.GenresProperty:               "Genres",
	vangogh_local_data.PropertiesProperty:           "Store Tags",
	vangogh_local_data.FeaturesProperty:             "Features",
	vangogh_local_data.LanguageCodeProperty:         "Language",
	vangogh_local_data.IncludesGamesProperty:        "Includes",
	vangogh_local_data.IsIncludedByGamesProperty:    "Included By",
	vangogh_local_data.RequiresGamesProperty:        "Requires",
	vangogh_local_data.IsRequiredByGamesProperty:    "Required By",
	vangogh_local_data.ProductTypeProperty:          "Product Type",
	vangogh_local_data.WishlistedProperty:           "Wishlisted",
	vangogh_local_data.OwnedProperty:                "Owned",
	vangogh_local_data.IsFreeProperty:               "Free",
	vangogh_local_data.IsDiscountedProperty:         "On Sale",
	vangogh_local_data.PreOrderProperty:             "Pre-order",
	vangogh_local_data.ComingSoonProperty:           "Coming Soon",
	vangogh_local_data.TBAProperty:                  "TBA",
	vangogh_local_data.InDevelopmentProperty:        "In Development",
	vangogh_local_data.IsUsingDOSBoxProperty:        "Using DOSBox",
	vangogh_local_data.IsUsingScummVMProperty:       "Using ScummVM",
	vangogh_local_data.TypesProperty:                "Data Type",
	vangogh_local_data.SteamReviewScoreDescProperty: "Steam Reviews",
	vangogh_local_data.SortProperty:                 "Sort",
	vangogh_local_data.DescendingProperty:           "Descending",
}

var gauginDigestibleProperties = []string{
	vangogh_local_data.TagIdProperty,
	vangogh_local_data.LocalTagsProperty,
	vangogh_local_data.LocalTagsProperty,
	vangogh_local_data.SteamTagsProperty,
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
	vangogh_local_data.SteamReviewScoreDescProperty,
}

var digestTitles = map[string]string{
	//vangogh_local_data.OperatingSystem
	"macos":   "macOS",
	"linux":   "Linux",
	"windows": "Windows",
	// vangogh_local_data.OwnedProperty, vangogh_local_data.WishlistedProperty, ...
	vangogh_local_data.TrueValue:  "Yes",
	vangogh_local_data.FalseValue: "No",
	//vangogh_local_data.SortProperty
	vangogh_local_data.GlobalReleaseDateProperty:  "Global Release Date",
	vangogh_local_data.GOGReleaseDateProperty:     "GOG.com Release Date",
	vangogh_local_data.GOGOrderDateProperty:       "GOG.com Order Date",
	vangogh_local_data.TitleProperty:              "Title",
	vangogh_local_data.RatingProperty:             "Rating",
	vangogh_local_data.DiscountPercentageProperty: "Discount Percentage",
	//vangogh_local_data.ProductTypeProperty
	vangogh_local_data.AccountProducts.String():  "Account Products",
	vangogh_local_data.ApiProductsV1.String():    "API Products V1",
	vangogh_local_data.ApiProductsV2.String():    "API Products V2",
	vangogh_local_data.Details.String():          "Account Product Details",
	vangogh_local_data.LicenceProducts.String():  "Licence Products",
	vangogh_local_data.Orders.String():           "Orders",
	vangogh_local_data.SteamAppNews.String():     "Steam App News",
	vangogh_local_data.SteamReviews.String():     "Steam Reviews",
	vangogh_local_data.SteamStorePage.String():   "Steam Store Page",
	vangogh_local_data.StoreProducts.String():    "Store Products",
	vangogh_local_data.WishlistProducts.String(): "Wishlist Products",
}

type searchProductsViewModel struct {
	Context          string
	Scope            string
	SearchProperties []string
	Query            map[string]string
	Digests          map[string][]string
	DigestsTitles    map[string]string
	Products         []listProductViewModel
	Limit            int
	Total            int
	Constrained      bool
	Path             string
}

func GetSearch(w http.ResponseWriter, r *http.Request) {

	// GET /search?(search_params)

	scope := ""
	for route, path := range searchRoutes {
		if r.URL.RawQuery != "" &&
			(strings.HasSuffix(path, r.URL.RawQuery) ||
				strings.HasSuffix(unconstrainedPath(path), r.URL.RawQuery)) {
			scope = route
			break
		}
	}

	q := r.URL.Query()

	spvm := &searchProductsViewModel{
		Scope:            scope,
		Context:          "filter-products",
		SearchProperties: gauginSearchProperties,
		Query:            make(map[string]string, len(q)),
		DigestsTitles:    digestTitles,
		Limit:            limit,
		Constrained:      !vangogh_local_data.FlagFromUrl(r.URL, "unconstrained"),
	}

	shortQuery := false
	queryProperties := append(gauginSearchProperties)
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

	if len(spvm.Query) > 0 {
		keys, err := getSearch(dc, q)
		if err != nil {
			http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
			return
		}

		spvm.Total = len(keys)
		if spvm.Total > spvm.Limit && spvm.Constrained {
			keys = keys[:limit]
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

		rdx, err := getRedux(dc, strings.Join(keys, ","), false, listReduxProperties...)

		if err != nil {
			http.Error(w, nod.ErrorStr("error getting all_redux"), http.StatusInternalServerError)
			return
		}

		lvm := listViewModelFromRedux(keys, rdx)
		spvm.Products = lvm.Products
	}

	// checking outside search action to account for empty query case
	if spvm.Total <= spvm.Limit {
		spvm.Constrained = false
	}

	digests, err := getDigests(dc, gauginDigestibleProperties...)

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

	if err != nil {
		http.Error(w, nod.ErrorStr("error getting digests"), http.StatusInternalServerError)
		return
	}
	spvm.Digests = digests

	gaugin_middleware.DefaultHeaders(w)

	if err := tmpl.ExecuteTemplate(w, "search-page", spvm); err != nil {
		http.Error(w, nod.ErrorStr("template error"), http.StatusInternalServerError)
		return
	}
}

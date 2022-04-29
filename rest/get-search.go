package rest

import (
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/http"
	"strings"
)

type searchQuery struct {
	Text        string
	Title       string
	Tags        string
	OS          string
	Developers  string
	Publisher   string
	Series      string
	Genres      string
	Properties  string
	Features    string
	Languages   string
	Includes    string
	IncludedBy  string
	Requires    string
	RequiredBy  string
	ProductType string
	Wishlisted  string
	Owned       string
	DataType    string
	Sort        string
	Desc        string
}

type searchProductsViewModel struct {
	Context  string
	Scope    string
	Query    searchQuery
	Digests  map[string][]string
	Products []listProductViewModel
}

func GetSearch(w http.ResponseWriter, r *http.Request) {
	dc := http.DefaultClient

	spvm := &searchProductsViewModel{
		Context: "filter-products",
		Scope:   r.URL.Query().Get("scope"),
	}

	q := r.URL.Query()
	spvm.Query = searchQuery{
		Text:        q.Get("text"),
		Title:       q.Get(vangogh_local_data.TitleProperty),
		Tags:        q.Get(vangogh_local_data.TagIdProperty),
		OS:          q.Get(vangogh_local_data.OperatingSystemsProperty),
		Developers:  q.Get(vangogh_local_data.DevelopersProperty),
		Publisher:   q.Get(vangogh_local_data.PublisherProperty),
		Series:      q.Get(vangogh_local_data.SeriesProperty),
		Genres:      q.Get(vangogh_local_data.GenresProperty),
		Properties:  q.Get(vangogh_local_data.PropertiesProperty),
		Features:    q.Get(vangogh_local_data.FeaturesProperty),
		Languages:   q.Get(vangogh_local_data.LanguageCodeProperty),
		Includes:    q.Get(vangogh_local_data.IncludesGamesProperty),
		IncludedBy:  q.Get(vangogh_local_data.IsIncludedByGamesProperty),
		Requires:    q.Get(vangogh_local_data.RequiresGamesProperty),
		RequiredBy:  q.Get(vangogh_local_data.IsRequiredByGamesProperty),
		ProductType: q.Get(vangogh_local_data.ProductTypeProperty),
		Wishlisted:  q.Get(vangogh_local_data.WishlistedProperty),
		Owned:       q.Get(vangogh_local_data.OwnedProperty),
		DataType:    q.Get(vangogh_local_data.TypesProperty),
		Sort:        q.Get("sort"),
		Desc:        q.Get("desc"),
	}

	emptyQuery := true
	for _, vs := range q {
		for _, v := range vs {
			if v != "" {
				emptyQuery = false
				break
			}
		}
	}

	if !emptyQuery {
		keys, err := getSearch(dc, q)
		if err != nil {
			http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
			return
		}

		rdx, err := getRedux(dc,
			strings.Join(keys, ","),
			vangogh_local_data.TitleProperty,
			vangogh_local_data.WishlistedProperty,
			vangogh_local_data.OwnedProperty,
			vangogh_local_data.DevelopersProperty,
			vangogh_local_data.PublisherProperty,
			vangogh_local_data.OperatingSystemsProperty,
			vangogh_local_data.TagIdProperty,
			vangogh_local_data.ProductTypeProperty)

		if err != nil {
			http.Error(w, nod.ErrorStr("error getting all_redux"), http.StatusInternalServerError)
			return
		}

		lvm := listViewModelFromRedux(keys, rdx)
		spvm.Products = lvm.Products
	}

	digests, err := getDigests(dc,
		vangogh_local_data.OperatingSystemsProperty,
		vangogh_local_data.GenresProperty,
		vangogh_local_data.PropertiesProperty,
		vangogh_local_data.FeaturesProperty,
		vangogh_local_data.LanguageCodeProperty,
		vangogh_local_data.ProductTypeProperty,
		vangogh_local_data.WishlistedProperty,
		vangogh_local_data.OwnedProperty,
		vangogh_local_data.TypesProperty)

	digests["sort"] = []string{vangogh_local_data.GOGReleaseDateProperty, vangogh_local_data.GOGOrderDateProperty, vangogh_local_data.TitleProperty}
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

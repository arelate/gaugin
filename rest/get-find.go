package rest

import (
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/http"
	"strings"
)

type searchQuery struct {
	Text       string
	Title      string
	Tags       string
	OS         string
	Developers string
	Publisher  string
	Series     string
	Genres     string
	Features   string
	Languages  string
	Includes   string
	IncludedBy string
	Requires   string
	RequiredBy string
}

type searchProductsViewModel struct {
	Context  string
	Query    searchQuery
	Products []listProductViewModel
}

func GetSearch(w http.ResponseWriter, r *http.Request) {
	dc := http.DefaultClient

	spvm := &searchProductsViewModel{
		Context: "search",
	}

	q := r.URL.Query()
	spvm.Query = searchQuery{
		Text:       q.Get("text"),
		Title:      q.Get("title"),
		Tags:       q.Get("tag"),
		OS:         q.Get("os"),
		Developers: q.Get("developers"),
		Publisher:  q.Get("publisher"),
		Series:     q.Get("series"),
		Genres:     q.Get("genres"),
		Features:   q.Get("features"),
		Languages:  q.Get("lang-code"),
		Includes:   q.Get("includes-games"),
		IncludedBy: q.Get("is-included-by-games"),
		Requires:   q.Get("requires-games"),
		RequiredBy: q.Get("is-required-by-games"),
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

	defaultHeaders(w)

	if err := tmpl.ExecuteTemplate(w, "search", spvm); err != nil {
		http.Error(w, nod.ErrorStr("template error"), http.StatusInternalServerError)
		return
	}
}

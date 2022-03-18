package api

import (
	"github.com/arelate/vangogh_local_data"
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

	keys, err := getSearch(dc, q)
	if err != nil {
		http.Error(w, "search error", http.StatusInternalServerError)
		return
	}

	redux, err := getRedux(dc,
		strings.Join(keys, ","),
		vangogh_local_data.TitleProperty,
		vangogh_local_data.DevelopersProperty,
		vangogh_local_data.PublisherProperty)
	if err != nil {
		http.Error(w, "error getting all_redux", http.StatusInternalServerError)
		return
	}

	lvm := listViewModelFromRedux(keys, redux)
	spvm.Products = lvm.Products

	defaultHeaders(w)

	if err := tmpl.ExecuteTemplate(w, "search", spvm); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}

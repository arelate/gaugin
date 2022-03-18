package api

import (
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_local_data"
	"net/http"
)

func getProductsList(
	pt vangogh_local_data.ProductType,
	mt gog_integration.Media,
	w http.ResponseWriter) {
	dc := http.DefaultClient

	keys, err := getKeys(dc, pt, mt)
	if err != nil {
		http.Error(w, "error getting keys", http.StatusInternalServerError)
		return
	}

	redux, err := getAllRedux(dc, pt, mt,
		vangogh_local_data.TitleProperty,
		vangogh_local_data.DevelopersProperty,
		vangogh_local_data.PublisherProperty)
	if err != nil {
		http.Error(w, "error getting all_redux", http.StatusInternalServerError)
		return
	}

	lvm := listViewModelFromRedux(keys, redux)
	lvm.Context = pt.String()

	defaultHeaders(w)

	if err := tmpl.ExecuteTemplate(w, "products-list", lvm); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}

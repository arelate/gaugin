package rest

import (
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/http"
)

func getProductsList(
	pt vangogh_local_data.ProductType,
	mt gog_integration.Media,
	w http.ResponseWriter) {
	dc := http.DefaultClient

	keys, err := getKeys(dc, pt, mt, -1)
	if err != nil {
		http.Error(w, nod.ErrorStr("error getting keys"), http.StatusInternalServerError)
		return
	}

	rdx, err := getAllRedux(dc, pt, mt,
		vangogh_local_data.TitleProperty,
		vangogh_local_data.DevelopersProperty,
		vangogh_local_data.PublisherProperty,
		vangogh_local_data.WishlistedProperty,
		vangogh_local_data.OwnedProperty,
		vangogh_local_data.OperatingSystemsProperty,
		vangogh_local_data.TagIdProperty,
		vangogh_local_data.ProductTypeProperty)
	if err != nil {
		http.Error(w, nod.ErrorStr("error getting all_redux"), http.StatusInternalServerError)
		return
	}

	lvm := listViewModelFromRedux(keys, rdx)
	lvm.Context = pt.String()

	defaultHeaders(w)

	if err := tmpl.ExecuteTemplate(w, "products-list", lvm); err != nil {
		http.Error(w, nod.ErrorStr("template exec error"), http.StatusInternalServerError)
		return
	}
}

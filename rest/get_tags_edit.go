package rest

import (
	"github.com/arelate/gaugin/gaugin_middleware"
	"github.com/arelate/gaugin/stencil_app"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/http"
)

func GetTagsEdit(w http.ResponseWriter, r *http.Request) {

	// GET /tags/edit?id

	id := r.URL.Query().Get("id")

	idRedux, _, err := getRedux(
		http.DefaultClient,
		id,
		false,
		vangogh_local_data.TitleProperty,
		vangogh_local_data.OwnedProperty,
		vangogh_local_data.TagIdProperty)

	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	digests, _, err := getDigests(http.DefaultClient, vangogh_local_data.TagIdProperty)

	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	selectedValues := make(map[string]bool)
	for _, v := range idRedux[id][vangogh_local_data.TagIdProperty] {
		selectedValues[v] = true
	}

	gaugin_middleware.DefaultHeaders(nil, w)

	if err := app.RenderPropertyEditor(
		id,
		idRedux[id][vangogh_local_data.TitleProperty][0],
		stencil_app.PropertyTitles[vangogh_local_data.TagIdProperty],
		idRedux[id][vangogh_local_data.OwnedProperty][0] == "true",
		"Account tags require product ownership",
		selectedValues,
		digests[vangogh_local_data.TagIdProperty],
		false,
		"/tags/apply",
		w); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}
}

package rest

import (
	"github.com/arelate/gaugin/gaugin_middleware"
	"github.com/arelate/gaugin/stencil_app"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/http"
)

func GetLocalTagsEdit(w http.ResponseWriter, r *http.Request) {

	// GET /local-tags/edit?id

	id := r.URL.Query().Get("id")

	idRedux, _, err := getRedux(
		http.DefaultClient,
		id,
		false,
		vangogh_local_data.TitleProperty,
		vangogh_local_data.LocalTagsProperty)

	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	digests, _, err := getDigests(http.DefaultClient, vangogh_local_data.LocalTagsProperty)

	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	gaugin_middleware.DefaultHeaders(nil, w)

	selectedValues := make(map[string]bool)
	for _, v := range idRedux[id][vangogh_local_data.LocalTagsProperty] {
		selectedValues[v] = true
	}

	if err := app.RenderPropertyEditor(
		id,
		idRedux[id][vangogh_local_data.TitleProperty][0],
		stencil_app.PropertyTitles[vangogh_local_data.LocalTagsProperty],
		true,
		"",
		selectedValues,
		digests[vangogh_local_data.LocalTagsProperty],
		true,
		"/local-tags/apply",
		w); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}
}

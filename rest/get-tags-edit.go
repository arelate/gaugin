package rest

import (
	"github.com/arelate/gaugin/gaugin_middleware"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/http"
)

func GetTagsEdit(w http.ResponseWriter, r *http.Request) {

	// GET /tags/edit?id

	id := r.URL.Query().Get("id")

	gaugin_middleware.DefaultHeaders(w)

	idRedux, err := getRedux(
		http.DefaultClient,
		id,
		false,
		vangogh_local_data.TitleProperty,
		vangogh_local_data.OwnedProperty,
		vangogh_local_data.TagIdProperty,
		vangogh_local_data.LocalTagsProperty)

	if err != nil {
		http.Error(w, nod.ErrorStr("error getting redux"), http.StatusInternalServerError)
		return
	}

	rdx, ok := idRedux[id]
	if !ok {
		http.Error(w, nod.ErrorStr("redux missing data for id %s", id), http.StatusInternalServerError)
		return
	}

	digests, err := getDigests(http.DefaultClient, vangogh_local_data.TagIdProperty, vangogh_local_data.LocalTagsProperty)

	if err != nil {
		http.Error(w, nod.ErrorStr("error getting redux"), http.StatusInternalServerError)
		return
	}

	tevm := &tagsEditViewModel{
		Context:      "tags",
		Id:           id,
		Title:        propertyFromRedux(rdx, vangogh_local_data.TitleProperty),
		Owned:        flagFromRedux(rdx, vangogh_local_data.OwnedProperty),
		AllTags:      propertiesFromRedux(digests, vangogh_local_data.TagIdProperty),
		AllLocalTags: propertiesFromRedux(digests, vangogh_local_data.LocalTagsProperty),
	}

	selectedTags := make(map[string]bool)
	for _, t := range propertiesFromRedux(rdx, vangogh_local_data.TagIdProperty) {
		selectedTags[t] = true
	}

	selectedLocalTags := make(map[string]bool)
	for _, t := range propertiesFromRedux(rdx, vangogh_local_data.LocalTagsProperty) {
		selectedLocalTags[t] = true
	}

	tevm.Tags = selectedTags
	tevm.LocalTags = selectedLocalTags

	if err := tmpl.ExecuteTemplate(w, "tags-edit-page", tevm); err != nil {
		http.Error(w, nod.ErrorStr("template exec error"), http.StatusInternalServerError)
		return
	}
}

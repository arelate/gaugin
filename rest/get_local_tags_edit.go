package rest

import (
	"github.com/arelate/gaugin/rest/compton_pages"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/kevlar"
	"github.com/boggydigital/nod"
	"net/http"
)

func GetLocalTagsEdit(w http.ResponseWriter, r *http.Request) {

	// GET /local-tags/edit?id

	id := r.URL.Query().Get("id")

	idRedux, err := getRedux(
		http.DefaultClient,
		id,
		false,
		vangogh_local_data.TitleProperty,
		vangogh_local_data.LocalTagsProperty,
		vangogh_local_data.ImageProperty,
		vangogh_local_data.DehydratedImageProperty)

	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	localTagsDigest, err := getDigests(http.DefaultClient, vangogh_local_data.LocalTagsProperty)

	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	rdx := kevlar.ReduxProxy(idRedux)

	selectedValues := make(map[string]any)
	if lt, ok := rdx.GetAllValues(vangogh_local_data.LocalTagsProperty, id); ok {
		for _, v := range lt {
			selectedValues[v] = nil
		}
	}

	localTags := make(map[string]string)
	for _, v := range localTagsDigest[vangogh_local_data.LocalTagsProperty] {
		localTags[v] = v
	}

	ltePage := compton_pages.TagsEditor(id, true, vangogh_local_data.LocalTagsProperty, localTags, selectedValues, rdx)
	if err := ltePage.WriteContent(w); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}
}

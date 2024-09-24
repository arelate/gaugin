package rest

import (
	"github.com/arelate/gaugin/rest/compton_pages"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/kevlar"
	"github.com/boggydigital/nod"
	"net/http"
)

func GetTagsEdit(w http.ResponseWriter, r *http.Request) {

	// GET /tags/edit?id

	id := r.URL.Query().Get("id")

	idRedux, err := getRedux(
		http.DefaultClient,
		id,
		false,
		vangogh_local_data.TitleProperty,
		vangogh_local_data.OwnedProperty,
		vangogh_local_data.TagIdProperty,
		vangogh_local_data.ImageProperty,
		vangogh_local_data.DehydratedImageProperty)

	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	rdx := kevlar.ReduxProxy(idRedux)

	tagNameRdx, err := getRedux(http.DefaultClient, "", true, vangogh_local_data.TagNameProperty)

	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	selectedValues := make(map[string]any)
	if tagIds, ok := rdx.GetAllValues(vangogh_local_data.TagIdProperty, id); ok {
		for _, v := range tagIds {
			selectedValues[v] = nil
		}
	}

	tagNames := make(map[string]string)

	for k, pv := range tagNameRdx {
		if v, ok := pv[vangogh_local_data.TagNameProperty]; ok && len(v) > 0 {
			tagNames[k] = v[0]
		}
	}

	owned := false
	if op, ok := rdx.GetLastVal(vangogh_local_data.OwnedProperty, id); ok && op == vangogh_local_data.TrueValue {
		owned = true
	}

	ltePage := compton_pages.TagsEditor(id, owned, vangogh_local_data.TagIdProperty, tagNames, selectedValues, rdx)
	if err := ltePage.WriteContent(w); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}
}

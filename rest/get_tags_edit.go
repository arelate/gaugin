package rest

import (
	"net/http"
	"strings"

	"github.com/arelate/gaugin/gaugin_middleware"
	"github.com/arelate/gaugin/view_models"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
)

func GetTagsEdit(w http.ResponseWriter, r *http.Request) {

	// GET /tags/edit?id

	id := r.URL.Query().Get("id")

	gaugin_middleware.DefaultHeaders(nil, w)

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

	digests, _, err := getDigests(http.DefaultClient, vangogh_local_data.TagIdProperty, vangogh_local_data.LocalTagsProperty)

	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	sb := &strings.Builder{}
	tevm := view_models.NewTagsEdit(id, idRedux[id], digests)

	if err := tmpl.ExecuteTemplate(sb, "tags-edit-form", tevm); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	if err := app.RenderPage(id, "Edit tags", sb.String(), w); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}
}

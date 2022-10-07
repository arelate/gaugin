package rest

import (
	"github.com/arelate/gaugin/gaugin_middleware"
	"net/http"
	"strings"

	"github.com/arelate/gaugin/view_models"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
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

	digests, _, err := getDigests(http.DefaultClient, vangogh_local_data.TagIdProperty, vangogh_local_data.LocalTagsProperty)

	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	sb := &strings.Builder{}
	tevm := view_models.NewTagsEdit(id, idRedux[id], digests)

	if err := tmpl.ExecuteTemplate(sb, "local-tags-edit-form", tevm); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	gaugin_middleware.DefaultHeaders(nil, w)

	if err := app.RenderPage(id, "Edit local tags", sb.String(), w); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}
}

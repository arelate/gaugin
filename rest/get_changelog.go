package rest

import (
	"github.com/arelate/gaugin/gaugin_middleware"
	"github.com/arelate/gaugin/stencil_app"
	"github.com/arelate/gaugin/view_models"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/http"
	"strings"
)

func GetChangelog(w http.ResponseWriter, r *http.Request) {

	// GET /changelog?id

	id := r.URL.Query().Get("id")

	idRedux, err := getRedux(
		http.DefaultClient,
		id,
		false,
		vangogh_local_data.ChangelogProperty)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	sb := &strings.Builder{}
	cvm := view_models.NewChangelog(idRedux[id])

	if err := tmpl.ExecuteTemplate(sb, "changelog-content", cvm); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	gaugin_middleware.DefaultHeaders(w)

	if err := app.RenderSection(id, stencil_app.ChangelogSection, sb.String(), w); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}
}

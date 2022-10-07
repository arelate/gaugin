package rest

import (
	"github.com/arelate/gaugin/gaugin_middleware"
	"github.com/arelate/gaugin/stencil_app"
	"net/http"
	"strings"

	"github.com/arelate/gaugin/view_models"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
)

func GetScreenshots(w http.ResponseWriter, r *http.Request) {

	// GET /screenshots?id

	id := r.URL.Query().Get("id")

	idRedux, _, err := getRedux(
		http.DefaultClient,
		id,
		false,
		vangogh_local_data.ScreenshotsProperty)

	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	sb := &strings.Builder{}
	svm := view_models.NewScreenshots(idRedux[id])

	if err := tmpl.ExecuteTemplate(sb, "screenshots-content", svm); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	gaugin_middleware.DefaultHeaders(nil, w)

	if err := app.RenderSection(id, stencil_app.ScreenshotsSection, sb.String(), w); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}
}

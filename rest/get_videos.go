package rest

import (
	"github.com/arelate/gaugin/gaugin_middleware"
	"github.com/arelate/gaugin/stencil_app"
	"github.com/boggydigital/kevlar"
	"net/http"
	"strings"

	"github.com/arelate/gaugin/view_models"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
)

func GetVideos(w http.ResponseWriter, r *http.Request) {

	// GET /videos?id

	id := r.URL.Query().Get("id")

	idRedux, err := getRedux(
		http.DefaultClient,
		id,
		false,
		vangogh_local_data.VideoIdProperty)

	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	mvRedux, err := getRedux(
		http.DefaultClient,
		"",
		true,
		vangogh_local_data.MissingVideoUrlProperty)

	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	rdx := kevlar.ReduxProxy(MergeIdPropertyValues(idRedux, mvRedux))

	sb := &strings.Builder{}
	vvm := view_models.NewVideos(id, rdx)

	if err := tmpl.ExecuteTemplate(sb, "videos-content", vvm); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	gaugin_middleware.DefaultHeaders(w)

	if err := app.RenderSection(id, stencil_app.VideosSection, sb.String(), w); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}
}

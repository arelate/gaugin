package rest

import (
	"github.com/arelate/gaugin/gaugin_middleware"
	"github.com/arelate/gaugin/view_models"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/http"
)

func GetVideos(w http.ResponseWriter, r *http.Request) {

	// GET /videos?id

	id := r.URL.Query().Get("id")

	gaugin_middleware.DefaultHeaders(w)

	idRedux, err := getRedux(
		http.DefaultClient,
		id,
		false,
		vangogh_local_data.VideoIdProperty)

	if err != nil {
		http.Error(w, nod.ErrorStr("error getting redux"), http.StatusInternalServerError)
		return
	}

	vvm := view_models.NewVideos(idRedux[id])

	if err := tmpl.ExecuteTemplate(w, "videos-page", vvm); err != nil {
		http.Error(w, nod.ErrorStr("template exec error"), http.StatusInternalServerError)
		return
	}
}

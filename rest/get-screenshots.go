package rest

import (
	"github.com/arelate/gaugin/gaugin_middleware"
	"github.com/arelate/gaugin/view_models"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/http"
)

func GetScreenshots(w http.ResponseWriter, r *http.Request) {

	// GET /screenshots?id

	id := r.URL.Query().Get("id")

	gaugin_middleware.DefaultHeaders(w)

	idRedux, err := getRedux(
		http.DefaultClient,
		id,
		false,
		vangogh_local_data.ScreenshotsProperty)

	if err != nil {
		http.Error(w, nod.ErrorStr("error getting redux"), http.StatusInternalServerError)
		return
	}

	svm := view_models.NewScreenshots(idRedux[id])

	if err := tmpl.ExecuteTemplate(w, "screenshots-page", svm); err != nil {
		http.Error(w, nod.ErrorStr("template exec error"), http.StatusInternalServerError)
		return
	}
}

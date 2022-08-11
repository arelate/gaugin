package rest

import (
	"github.com/arelate/gaugin/gaugin_middleware"
	"github.com/arelate/gaugin/view_models"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/http"
)

func GetDescription(w http.ResponseWriter, r *http.Request) {

	// GET /description?id

	id := r.URL.Query().Get("id")

	idRedux, err := getRedux(
		http.DefaultClient,
		id,
		false,
		vangogh_local_data.DescriptionOverviewProperty,
		vangogh_local_data.DescriptionFeaturesProperty,
		vangogh_local_data.AdditionalRequirementsProperty,
		vangogh_local_data.CopyrightsProperty)

	if err != nil {
		http.Error(w, nod.ErrorStr("error getting redux"), http.StatusInternalServerError)
		return
	}

	gaugin_middleware.DefaultHeaders(w)

	dvm := view_models.NewDescription(idRedux[id])

	if err := tmpl.ExecuteTemplate(w, "description-page", dvm); err != nil {
		http.Error(w, nod.ErrorStr("template exec error"), http.StatusInternalServerError)
		return
	}
}

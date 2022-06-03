package rest

import (
	"github.com/arelate/gaugin/gaugin_middleware"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"html/template"
	"net/http"
)

func GetDescription(w http.ResponseWriter, r *http.Request) {

	// GET /description?id

	id := r.URL.Query().Get("id")

	idRedux, err := getRedux(
		http.DefaultClient,
		id,
		vangogh_local_data.DescriptionOverviewProperty,
		vangogh_local_data.DescriptionFeaturesProperty,
		vangogh_local_data.AdditionalRequirementsProperty,
		vangogh_local_data.CopyrightsProperty)

	if err != nil {
		http.Error(w, nod.ErrorStr("error getting redux"), http.StatusInternalServerError)
		return
	}

	gaugin_middleware.DefaultHeaders(w)

	dvm := &descriptionViewModel{Context: "description"}

	for _, rdx := range idRedux {

		//Description content preparation includes the following steps:
		//1) combining DescriptionOverview and DescriptionFeatures
		//2) replacing implicit list in DescriptionFeatures with explicit list
		//3) rewriting https://items.gog.com/... links to gaugin
		//4) rewriting https://www.gog.com/game/... links to gaugin

		desc := propertyFromRedux(rdx, vangogh_local_data.DescriptionOverviewProperty)
		desc += implicitToExplicitList(propertyFromRedux(rdx, vangogh_local_data.DescriptionFeaturesProperty))

		desc = rewriteDescriptionItemsLinks(desc)
		desc = rewriteDescriptionGameLinks(desc)

		dvm.Description = template.HTML(desc)

		if ar := propertyFromRedux(rdx, vangogh_local_data.AdditionalRequirementsProperty); ar != "" {
			dvm.AdditionalRequirements = template.HTML(ar)
		}

		if c := propertyFromRedux(rdx, vangogh_local_data.CopyrightsProperty); c != "" {
			dvm.Copyrights = template.HTML(c)
		}
	}

	if err := tmpl.ExecuteTemplate(w, "description-page", dvm); err != nil {
		http.Error(w, nod.ErrorStr("template exec error"), http.StatusInternalServerError)
		return
	}

}

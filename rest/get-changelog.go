package rest

import (
	"github.com/arelate/gaugin/gaugin_middleware"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"html/template"
	"net/http"
)

func GetChangelog(w http.ResponseWriter, r *http.Request) {

	// GET /changelog?id

	id := r.URL.Query().Get("id")

	idRedux, err := getRedux(
		http.DefaultClient,
		id,
		vangogh_local_data.ChangelogProperty)

	if err != nil {
		http.Error(w, nod.ErrorStr("error getting redux"), http.StatusInternalServerError)
		return
	}

	gaugin_middleware.DefaultHeaders(w)

	cvm := &changelogViewModel{
		Context:   "changelog",
		Changelog: template.HTML(propertyFromRedux(idRedux[id], vangogh_local_data.ChangelogProperty)),
	}

	if err := tmpl.ExecuteTemplate(w, "changelog-page", cvm); err != nil {
		http.Error(w, nod.ErrorStr("template exec error"), http.StatusInternalServerError)
		return
	}
}

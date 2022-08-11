package rest

import (
	"github.com/arelate/gaugin/gaugin_middleware"
	"github.com/arelate/gaugin/view_models"
	"github.com/boggydigital/nod"
	"net/http"
)

func GetDownloads(w http.ResponseWriter, r *http.Request) {

	// GET /downloads?id

	id := r.URL.Query().Get("id")

	gaugin_middleware.DefaultHeaders(w)

	dvm, err := getDownloadsViewModel(id, r.Header.Get("User-Agent"))
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "downloads-page", dvm); err != nil {
		http.Error(w, nod.ErrorStr("template exec error"), http.StatusInternalServerError)
		return
	}
}

func getDownloadsViewModel(id string, userAgent string) (*view_models.Downloads, error) {

	//we specifically get /downloads and not /data&product-type=details because of Details
	//format complexities, see gog_integration/details.go/GetGameDownloads comment
	dls, err := getDownloads(http.DefaultClient, id, operatingSystems, languageCodes)
	if err != nil {
		return nil, err
	}

	dvm := view_models.NewDownloads(userAgent, dls)

	return dvm, nil
}

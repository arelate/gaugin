package rest

import (
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"net/http"
)

func GetFiles(w http.ResponseWriter, r *http.Request) {

	// GET /files?manual-url

	q := r.URL.Query()

	manualUrl := q.Get("manual-url")

	if manualUrl != "" {

		idRedux, _, err := getRedux(http.DefaultClient, manualUrl, false, vangogh_local_data.LocalManualUrlProperty)
		if err != nil {
			http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
			return
		}

		rxa := kvas.NewIRAProxy(idRedux)

		relLocalFilePath, ok := rxa.GetFirstVal(vangogh_local_data.LocalManualUrlProperty, manualUrl)
		if !ok {
			http.Error(w, nod.ErrorStr("no file for manual-url %s", manualUrl), http.StatusNotFound)
			return
		}

		http.Redirect(w, r, "/local-file/"+relLocalFilePath, http.StatusPermanentRedirect)
	} else {
		http.Error(w, nod.ErrorStr("missing manual-url"), http.StatusNotFound)
	}
}

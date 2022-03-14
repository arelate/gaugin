package api

import (
	"fmt"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/http"
)

func GetFiles(w http.ResponseWriter, r *http.Request) {

	// GET /files?manual-url

	if r.Method != http.MethodGet {
		err := fmt.Errorf("unsupported method")
		http.Error(w, nod.Error(err).Error(), 405)
		return
	}

	q := r.URL.Query()

	manualUrl := q.Get("manual-url")

	if manualUrl != "" {
		relLocalFilePath, ok := rxa.GetFirstVal(vangogh_local_data.LocalManualUrlProperty, manualUrl)
		if !ok {
			http.Error(w, fmt.Sprintf("no file for manual-url %s", manualUrl), http.StatusNotFound)
			return
		}

		http.Redirect(w, r, "/local-file/"+relLocalFilePath, http.StatusPermanentRedirect)
	} else {
		http.Error(w, "missing manual-url", http.StatusNotFound)
	}
}

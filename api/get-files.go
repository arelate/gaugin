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
	if manualUrl == "" {
		err := fmt.Errorf("empty manual-url")
		http.Error(w, nod.Error(err).Error(), 400)
		return
	}

	relLocalFilePath, ok := rxa.GetFirstVal(vangogh_local_data.LocalManualUrlProperty, manualUrl)
	if !ok {
		http.Error(w, fmt.Sprintf("no file for manual-url %s", manualUrl), http.StatusNotFound)
		return
	}

	if absLocalFilePath := vangogh_local_data.AbsDownloadDirFromRel(relLocalFilePath); absLocalFilePath != "" {
		http.ServeFile(w, r, absLocalFilePath)
	} else {
		http.Error(w, fmt.Sprintf("no local file for manual-url %s", manualUrl), http.StatusNotFound)
		return
	}
}

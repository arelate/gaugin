package api

import (
	"fmt"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/http"
)

var (
	operatingSystems []vangogh_local_data.OperatingSystem
	languageCodes    []string
)

func SetDownloadsOperatingSystems(os []vangogh_local_data.OperatingSystem) {
	operatingSystems = os
}

func SetDownloadsLanguageCodes(lc []string) {
	languageCodes = lc
}

func GetDownloads(w http.ResponseWriter, r *http.Request) {

	// GET /downloads?id

	if r.Method != http.MethodGet {
		err := fmt.Errorf("unsupported method")
		http.Error(w, nod.Error(err).Error(), 405)
		return
	}

	dc := http.DefaultClient
	q := r.URL.Query()
	id := q.Get("id")

	dl, err := getDownloads(dc, id, operatingSystems, languageCodes)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "product", dl); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}

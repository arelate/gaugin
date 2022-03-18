package api

import (
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/http"
	"strings"
)

func GetProduct(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")

	idRedux, err := getRedux(http.DefaultClient, id, vangogh_local_data.ReduxProperties()...)
	if err != nil {
		http.Error(w, "error getting redux", http.StatusInternalServerError)
		return
	}

	defaultHeaders(w)

	pvm, err := productViewModelFromRedux(idRedux)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	validVideos := make([]string, 0, len(pvm.Videos))
	for _, v := range pvm.Videos {
		if strings.Contains(v, "(") {
			continue
		}
		validVideos = append(validVideos, v)
	}
	pvm.Videos = validVideos

	dl, err := getDownloads(http.DefaultClient, id, operatingSystems, languageCodes)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	pvm.Downloads = dl

	if err := tmpl.ExecuteTemplate(w, "product", pvm); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}

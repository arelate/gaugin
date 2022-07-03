package rest

import (
	"github.com/arelate/gaugin/gaugin_middleware"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/http"
	"strings"
)

func GetVideos(w http.ResponseWriter, r *http.Request) {

	// GET /videos?id

	id := r.URL.Query().Get("id")

	gaugin_middleware.DefaultHeaders(w)

	idRedux, err := getRedux(
		http.DefaultClient,
		id,
		vangogh_local_data.VideoIdProperty)

	if err != nil {
		http.Error(w, nod.ErrorStr("error getting redux"), http.StatusInternalServerError)
		return
	}

	vvm := &videosViewModel{
		Context:      "iframe",
		LocalVideos:  make([]string, 0),
		RemoteVideos: make([]string, 0),
	}

	// filter videos to distinguish between locally available and remote videos

	for _, v := range propertiesFromRedux(idRedux[id], vangogh_local_data.VideoIdProperty) {
		if !strings.Contains(v, "(") {
			vvm.LocalVideos = append(vvm.LocalVideos, v)
		} else {
			vvm.RemoteVideos = append(vvm.RemoteVideos, v)
		}
	}

	if err := tmpl.ExecuteTemplate(w, "videos-page", vvm); err != nil {
		http.Error(w, nod.ErrorStr("template exec error"), http.StatusInternalServerError)
		return
	}
}

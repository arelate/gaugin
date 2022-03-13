package api

import "net/http"

func GetVideos(w http.ResponseWriter, r *http.Request) {
	vu := videoUrl(r.URL.Query().Get("id"))
	getMedia(vu, w)
}

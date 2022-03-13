package api

import (
	"net/http"
)

func GetImages(w http.ResponseWriter, r *http.Request) {
	iu := imageUrl(r.URL.Query().Get("id"))
	getMedia(iu, w)
}

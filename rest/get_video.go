package rest

import (
	"fmt"
	"net/http"

	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
)

func GetVideo(w http.ResponseWriter, r *http.Request) {

	// GET /video?id

	q := r.URL.Query()
	videoId := q.Get("id")
	if videoId == "" {
		err := fmt.Errorf("empty video id")
		http.Error(w, nod.Error(err).Error(), http.StatusBadRequest)
		return
	}
	if localVideoPath, err := vangogh_local_data.AbsLocalVideoPath(videoId); err == nil && localVideoPath != "" {
		w.Header().Set("Cache-Control", "max-age=31536000")
		http.ServeFile(w, r, localVideoPath)
	} else {
		if err == nil {
			err = fmt.Errorf("no local video for id %s", videoId)
		}
		http.Error(w, nod.Error(err).Error(), http.StatusNotFound)
	}
}

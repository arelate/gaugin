package api

import (
	"fmt"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/http"
)

func GetVideos(w http.ResponseWriter, r *http.Request) {

	// GET /videos?id

	if r.Method != http.MethodGet {
		err := fmt.Errorf("unsupported method")
		http.Error(w, nod.Error(err).Error(), 405)
		return
	}

	q := r.URL.Query()
	videoId := q.Get("id")
	if videoId == "" {
		err := fmt.Errorf("empty video id")
		http.Error(w, nod.Error(err).Error(), 400)
		return
	}
	if localVideoPath := vangogh_local_data.AbsLocalVideoPath(videoId); localVideoPath != "" {
		w.Header().Set("Cache-Control", "max-age=31536000")
		http.ServeFile(w, r, localVideoPath)
	} else {
		_ = nod.Error(fmt.Errorf("no local video for id %s", videoId))
		http.NotFound(w, r)
	}
}

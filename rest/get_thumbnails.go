package rest

import (
	"fmt"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/http"
)

func GetThumbnails(w http.ResponseWriter, r *http.Request) {

	// GET /thumbnails?id

	q := r.URL.Query()
	videoId := q.Get("id")
	if videoId == "" {
		err := fmt.Errorf("empty video id")
		http.Error(w, nod.Error(err).Error(), http.StatusBadRequest)
		return
	}
	if localThumbnailPath := vangogh_local_data.AbsLocalVideoThumbnailPath(videoId); localThumbnailPath != "" {
		w.Header().Set("Cache-Control", "max-age=31536000")
		http.ServeFile(w, r, localThumbnailPath)
	} else {
		_ = nod.Error(fmt.Errorf("no local thumbnail for id %s", videoId))
		http.NotFound(w, r)
	}
}

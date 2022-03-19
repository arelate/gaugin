package rest

import (
	"fmt"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/http"
)

func GetImages(w http.ResponseWriter, r *http.Request) {

	// GET /images?id

	if r.Method != http.MethodGet {
		err := fmt.Errorf("unsupported method")
		http.Error(w, nod.Error(err).Error(), 405)
		return
	}

	q := r.URL.Query()
	imageId := q.Get("id")
	if imageId == "" {
		err := fmt.Errorf("empty image id")
		http.Error(w, nod.Error(err).Error(), 400)
		return
	}
	if localImagePath := vangogh_local_data.AbsLocalImagePath(imageId); localImagePath != "" {
		w.Header().Set("Cache-Control", "max-age=31536000")
		http.ServeFile(w, r, localImagePath)
	} else {
		_ = nod.Error(fmt.Errorf("no local image for id %s", imageId))
		http.NotFound(w, r)
	}
}

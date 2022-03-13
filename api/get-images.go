package api

import (
	"fmt"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/http"
)

//func GetImages(w http.ResponseWriter, r *http.Request) {
//	iu := imageUrl(r.URL.Query().Get("id"))
//	dc := http.DefaultClient
//
//	resp, err := dc.Get(iu.String())
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	defer resp.Body.Close()
//
//	w.Header().Add("Content-Type", resp.Header.Get("Content-Type"))
//	w.Header().Add("Last-Modified", resp.Header.Get("Last-Modified"))
//
//	w.WriteHeader(http.StatusOK)
//
//	if _, err := io.Copy(w, resp.Body); err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//}

func GetImages(w http.ResponseWriter, r *http.Request) {

	// GET /v1/images?id

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
		http.ServeFile(w, r, localImagePath)
	} else {
		_ = nod.Error(fmt.Errorf("no local image for id %s", imageId))
		http.NotFound(w, r)
	}
}

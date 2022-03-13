package api

import (
	"fmt"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/http"
)

//const videoCacheDir = "/tmp"
//
//func GetVideos(w http.ResponseWriter, r *http.Request) {
//	id := r.URL.Query().Get("id")
//
//	dc := http.DefaultClient
//	vu := videoUrl(id)
//
//	resp, err := dc.Get(vu.String())
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	defer resp.Body.Close()
//
//	cachedVideoPath := filepath.Join(videoCacheDir, id)
//	if _, err := os.Stat(cachedVideoPath); os.IsNotExist(err) {
//		cachedVideoFile, err := os.Create(cachedVideoPath)
//		if err != nil {
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//			return
//		}
//
//		if _, err := io.Copy(cachedVideoFile, resp.Body); err != nil {
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//			return
//		}
//
//		mod, err := time.Parse(time.RFC1123, resp.Header.Get("Last-Modified"))
//		if err != nil {
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//			return
//		}
//
//		os.Chtimes(cachedVideoPath, mod, mod)
//	}
//
//	http.ServeFile(w, r, cachedVideoPath)
//}

func GetVideos(w http.ResponseWriter, r *http.Request) {

	// GET /v1/videos?id

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
		http.ServeFile(w, r, localVideoPath)
	} else {
		_ = nod.Error(fmt.Errorf("no local video for id %s", videoId))
		http.NotFound(w, r)
	}
}

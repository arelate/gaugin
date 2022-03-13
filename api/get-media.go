package api

import (
	"io"
	"log"
	"net/http"
	"net/url"
)

func getMedia(url *url.URL, w http.ResponseWriter) {
	dc := http.DefaultClient

	resp, err := dc.Get(url.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.Header().Add("Content-Type", resp.Header.Get("Content-Type"))
	w.Header().Add("Last-Modified", resp.Header.Get("Last-Modified"))

	w.WriteHeader(http.StatusOK)

	if _, err := io.Copy(w, resp.Body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

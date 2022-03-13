package api

import (
	"io"
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

	if _, err := io.Copy(w, resp.Body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

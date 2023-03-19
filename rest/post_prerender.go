package rest

import (
	"bytes"
	"fmt"
	"github.com/boggydigital/nod"
	"io"
	"net/http"
)

var staticContent map[string][]byte

func getStaticContent(w http.ResponseWriter, r *http.Request) bool {
	key := r.URL.Path
	if r.URL.RawQuery != "" {
		key += "?" + r.URL.RawQuery
	}
	if bs, ok := staticContent[key]; ok {
		br := bytes.NewReader(bs)
		if _, err := io.Copy(w, br); err != nil {
			http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
			return false
		}
		return true
	}
	return false
}

func PostPrerender(w http.ResponseWriter, r *http.Request) {

	// POST /prerender

	if staticContent == nil {
		staticContent = make(map[string][]byte)
	}

	paths := []string{
		"/updates",
	}

	for _, p := range searchRoutes() {
		paths = append(paths, p)
	}

	host := fmt.Sprintf("http://localhost:%d", port)

	for _, p := range paths {
		resp, err := http.DefaultClient.Get(host + p)
		if err != nil {
			http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		bs := make([]byte, 0, 1024*1024)
		bb := bytes.NewBuffer(bs)

		if _, err := io.Copy(bb, resp.Body); err != nil {
			http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
			return
		}

		staticContent[p] = bb.Bytes()
	}

	if _, err := io.WriteString(w, "ok"); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

}

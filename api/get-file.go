package api

import (
	"net/http"
	"os"
	"path/filepath"
)

func GetFile(w http.ResponseWriter, r *http.Request) {
	p, err := filepath.Rel("/", r.URL.Path)
	if err != nil {
		http.Error(w, "invalid path", http.StatusRequestedRangeNotSatisfiable)
		return
	}
	if _, err := os.Stat(p); err == nil {
		http.ServeFile(w, r, p)
	} else {
		http.Error(w, "file not found", http.StatusNotFound)
	}
}

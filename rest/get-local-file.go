package rest

import (
	"fmt"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/http"
	"path/filepath"
)

func GetLocalFile(w http.ResponseWriter, r *http.Request) {

	// GET /local-file/{rel-local-path}

	localPath, err := filepath.Rel("/local-file/", r.URL.Path)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusMisdirectedRequest)
		return
	}

	if absLocalFilePath := vangogh_local_data.AbsDownloadDirFromRel(localPath); absLocalFilePath != "" {
		http.ServeFile(w, r, absLocalFilePath)
	} else {
		_ = nod.Error(fmt.Errorf("file %s not found", absLocalFilePath))
		http.NotFound(w, r)
	}
}

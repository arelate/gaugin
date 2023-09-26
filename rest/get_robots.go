package rest

import (
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/stencil/stencil_rest"
	"net/http"
)

func GetRobotsTxt(w http.ResponseWriter, r *http.Request) {
	aifp, err := vangogh_local_data.GetAbsDir(vangogh_local_data.InputFiles)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusNotFound)
		return
	}

	stencil_rest.GetRobotsTxt(aifp, w, r)
}

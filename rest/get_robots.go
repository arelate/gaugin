package rest

import (
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/pasu"
	"github.com/boggydigital/stencil/stencil_rest"
	"net/http"
)

func GetRobotsTxt(w http.ResponseWriter, r *http.Request) {
	// BUG: this needs directories rewrite
	aifp, err := pasu.GetAbsDir(vangogh_local_data.Input)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusNotFound)
		return
	}

	stencil_rest.GetRobotsTxt(aifp, w, r)
}

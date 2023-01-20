package rest

import (
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/stencil/stencil_rest"
	"net/http"
)

func GetRobotsTxt(w http.ResponseWriter, r *http.Request) {
	stencil_rest.GetRobotsTxt(vangogh_local_data.Pwd(), w, r)
}

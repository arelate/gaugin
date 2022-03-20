package rest

import (
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_local_data"
	"net/http"
)

func GetDownloads(w http.ResponseWriter, r *http.Request) {
	getProductsList(vangogh_local_data.Details, gog_integration.Game, w)
}

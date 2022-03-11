package api

import (
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_local_data"
	"net/http"
)

func GetStore(w http.ResponseWriter, r *http.Request) {
	getProductsList(vangogh_local_data.StoreProducts, gog_integration.Game, w)
}

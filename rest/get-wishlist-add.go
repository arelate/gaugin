package rest

import (
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/http"
)

func GetWishlistAdd(w http.ResponseWriter, r *http.Request) {

	// GET /wishlist/add?id

	id := r.URL.Query().Get(vangogh_local_data.IdProperty)

	if err := putWishlist(http.DefaultClient, id); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/product?id="+id, http.StatusTemporaryRedirect)
}

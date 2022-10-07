package rest

import (
	"github.com/boggydigital/nod"
	"net/http"
)

func GetTagsApply(w http.ResponseWriter, r *http.Request) {

	// GET /tags/apply

	if err := r.ParseForm(); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusBadRequest)
		return
	}

	var id string
	if len(r.Form["id"]) > 0 {
		id = r.Form["id"][0]
	}

	owned := false
	if len(r.Form["owned"]) > 0 {
		owned = r.Form["owned"][0] == "true"
	}

	if owned {
		//don't skip if tags are empty as this might be a signal to remove existing tags
		tags := r.Form["tag"]
		if err := patchTag(http.DefaultClient, id, tags); err != nil {
			http.Error(w, nod.Error(err).Error(), http.StatusBadRequest)
			return
		}
	}

	http.Redirect(w, r, "/product?id="+id, http.StatusTemporaryRedirect)
}

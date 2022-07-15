package rest

import (
	"github.com/boggydigital/nod"
	"net/http"
)

func PostTagsApply(w http.ResponseWriter, r *http.Request) {

	// POST /tags/apply

	if err := r.ParseForm(); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusBadRequest)
		return
	}

	var id string
	if len(r.Form["id"]) > 0 {
		id = r.Form["id"][0]
	}

	tags := r.Form["tag"]
	if err := patchTag(http.DefaultClient, id, tags); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusBadRequest)
		return
	}

	localTags := r.Form["local-tag"]
	if err := patchLocalTag(http.DefaultClient, id, localTags); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/product?id="+id, http.StatusTemporaryRedirect)
}

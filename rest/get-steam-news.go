package rest

import (
	"github.com/arelate/gaugin/gaugin_middleware"
	"github.com/arelate/gaugin/view_models"
	"github.com/boggydigital/nod"
	"net/http"
)

func GetSteamNews(w http.ResponseWriter, r *http.Request) {

	// GET /steam-news?id

	id := r.URL.Query().Get("id")

	gaugin_middleware.DefaultHeaders(w)

	san, err := getSteamAppNews(http.DefaultClient, id)
	if err != nil {
		http.Error(w, nod.ErrorStr("error getting steam news"), http.StatusInternalServerError)
		return
	}

	sanvm := view_models.NewSteamNews(san)

	if err := tmpl.ExecuteTemplate(w, "steam-news-page", sanvm); err != nil {
		http.Error(w, nod.ErrorStr("template exec error"), http.StatusInternalServerError)
		return
	}
}

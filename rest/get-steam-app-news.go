package rest

import (
	"github.com/arelate/gaugin/gaugin_middleware"
	"github.com/boggydigital/nod"
	"net/http"
)

func GetSteamAppNews(w http.ResponseWriter, r *http.Request) {

	// GET /steam-app-news?id

	id := r.URL.Query().Get("id")

	gaugin_middleware.DefaultHeaders(w)

	sanvm := &steamAppNewsViewModel{Context: "steam-app-news"}

	var err error
	sanvm.SteamAppNews, err = getSteamAppNews(http.DefaultClient, id)
	if err != nil {
		http.Error(w, nod.ErrorStr("error getting steam app news"), http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "steam-app-news-page", sanvm); err != nil {
		http.Error(w, nod.ErrorStr("template exec error"), http.StatusInternalServerError)
		return
	}
}

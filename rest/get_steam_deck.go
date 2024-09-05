package rest

import (
	"github.com/arelate/gaugin/gaugin_middleware"
	"github.com/arelate/gaugin/stencil_app"
	"github.com/arelate/gaugin/view_models"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/http"
	"strings"
)

func GetSteamDeck(w http.ResponseWriter, r *http.Request) {

	// GET /steam-deck?id

	id := r.URL.Query().Get("id")

	sdacr, err := getSteamDeckReport(http.DefaultClient, id)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	titleResult, err := getRedux(http.DefaultClient, id, false, vangogh_local_data.TitleProperty)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	title := ""
	if tt, ok := titleResult[id][vangogh_local_data.TitleProperty]; ok && len(tt) > 0 {
		title = tt[0]
	}

	sb := &strings.Builder{}
	sdvm := view_models.NewSteamDeck(title, sdacr)

	if err := tmpl.ExecuteTemplate(sb, "steam-deck-content", sdvm); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	gaugin_middleware.DefaultHeaders(w)

	if err := app.RenderSection(id, stencil_app.SteamDeckSection, sb.String(), w); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}
}

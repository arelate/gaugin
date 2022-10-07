package rest

import (
	"github.com/arelate/gaugin/gaugin_middleware"
	"github.com/arelate/gaugin/stencil_app"
	"github.com/arelate/gaugin/view_models"
	"github.com/boggydigital/nod"
	"net/http"
	"strings"
)

func GetSteamReviews(w http.ResponseWriter, r *http.Request) {

	// GET /steam-reviews?id

	id := r.URL.Query().Get("id")

	sar, _, err := getSteamReviews(http.DefaultClient, id)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	sb := &strings.Builder{}
	srvm := view_models.NewSteamReviews(sar)

	if err := tmpl.ExecuteTemplate(sb, "steam-reviews-content", srvm); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	gaugin_middleware.DefaultHeaders(nil, w)

	if err := app.RenderSection(id, stencil_app.SteamReviewsSection, sb.String(), w); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}
}

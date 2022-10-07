package rest

import (
	"github.com/arelate/gaugin/stencil_app"
	"net/http"
	"strings"
	"time"

	"github.com/arelate/gaugin/gaugin_middleware"
	"github.com/arelate/gaugin/view_models"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
)

func GetChangelog(w http.ResponseWriter, r *http.Request) {

	// GET /changelog?id

	id := r.URL.Query().Get("id")

	st := gaugin_middleware.NewServerTimings()
	start := time.Now()

	idRedux, cached, err := getRedux(
		http.DefaultClient,
		id,
		false,
		vangogh_local_data.ChangelogProperty)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	if cached {
		st.SetFlag("getRedux-cached")
	}
	st.Set("getRedux", time.Since(start).Milliseconds())

	sb := &strings.Builder{}
	cvm := view_models.NewChangelog(idRedux[id])

	if err := tmpl.ExecuteTemplate(sb, "changelog-content", cvm); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	gaugin_middleware.DefaultHeaders(st, w)

	if err := app.RenderSection(id, stencil_app.ChangelogSection, sb.String(), w); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}
}

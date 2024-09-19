package rest

import (
	"github.com/arelate/gaugin/rest/compton_data"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/compton/consts/color"
	"github.com/boggydigital/compton/consts/direction"
	"github.com/boggydigital/compton/elements/els"
	"github.com/boggydigital/compton/elements/flex_items"
	"github.com/boggydigital/compton/elements/fspan"
	"github.com/boggydigital/compton/elements/iframe_expand"
	"github.com/boggydigital/compton/elements/recipes"
	"github.com/boggydigital/nod"
	"net/http"
)

func GetChangelog(w http.ResponseWriter, r *http.Request) {

	// GET /changelog?id

	id := r.URL.Query().Get("id")

	idRedux, err := getRedux(
		http.DefaultClient,
		id,
		false,
		vangogh_local_data.ChangelogProperty)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	var changelog []string
	if rdx := idRedux[id]; rdx != nil {
		changelog = rdx[vangogh_local_data.ChangelogProperty]
	}

	section := compton_data.ChangelogSection
	ifc := iframe_expand.IframeExpandContent(section, compton_data.SectionTitles[section])

	pageStack := flex_items.FlexItems(ifc, direction.Column)
	ifc.Append(pageStack)

	if len(changelog) == 0 {
		fs := fspan.Text(ifc, "Changelog is not available for this product").
			ForegroundColor(color.Subtle)
		pageStack.Append(recipes.Center(ifc, fs))
	}

	for _, log := range changelog {
		pageStack.Append(els.Text(log))
	}

	if err := ifc.WriteContent(w); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
	}
}

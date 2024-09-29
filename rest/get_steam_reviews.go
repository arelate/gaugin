package rest

import (
	"github.com/arelate/gaugin/rest/compton_data"
	"github.com/arelate/gaugin/rest/compton_fragments"
	"github.com/arelate/gaugin/rest/gaugin_styles"
	"github.com/boggydigital/compton/consts/color"
	"github.com/boggydigital/compton/consts/direction"
	"github.com/boggydigital/compton/elements/els"
	"github.com/boggydigital/compton/elements/flex_items"
	"github.com/boggydigital/compton/elements/fspan"
	"github.com/boggydigital/compton/elements/iframe_expand"
	"github.com/boggydigital/nod"
	"net/http"
)

func GetSteamReviews(w http.ResponseWriter, r *http.Request) {

	// GET /steam-reviews?id

	id := r.URL.Query().Get("id")

	sar, err := getSteamReviews(http.DefaultClient, id)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	section := compton_data.SteamReviewsSection
	ifc := iframe_expand.IframeExpandContent(section, compton_data.SectionTitles[section]).
		AppendStyle(gaugin_styles.SteamReviews)

	pageStack := flex_items.FlexItems(ifc, direction.Column)
	ifc.Append(pageStack)

	if len(sar.Reviews) == 0 {
		fs := fspan.Text(ifc, "Steam reviews are not available for this product").
			ForegroundColor(color.Gray)
		pageStack.Append(flex_items.Center(ifc, fs))
	}

	for ii, review := range sar.Reviews {
		if srf := compton_fragments.SteamReview(ifc, review); srf != nil {
			pageStack.Append(srf)
		}
		if ii < len(sar.Reviews)-1 {
			pageStack.Append(els.Hr())
		}
	}

	if err := ifc.WriteContent(w); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
	}
}

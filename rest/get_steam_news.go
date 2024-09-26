package rest

import (
	"github.com/arelate/gaugin/rest/compton_data"
	"github.com/arelate/gaugin/rest/gaugin_styles"
	"github.com/arelate/southern_light/steam_integration"
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/direction"
	"github.com/boggydigital/compton/consts/weight"
	"github.com/boggydigital/compton/elements/details_summary"
	"github.com/boggydigital/compton/elements/flex_items"
	"github.com/boggydigital/compton/elements/fspan"
	"github.com/boggydigital/compton/elements/iframe_expand"
	"github.com/boggydigital/nod"
	"net/http"
)

func GetSteamNews(w http.ResponseWriter, r *http.Request) {

	// GET /steam-news?id

	id := r.URL.Query().Get("id")

	san, err := getSteamNews(http.DefaultClient, id)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	section := compton_data.SteamNewsSection
	ifc := iframe_expand.IframeExpandContent(section, compton_data.SectionTitles[section]).
		AppendStyle(gaugin_styles.SteamReviews)

	pageStack := flex_items.FlexItems(ifc, direction.Column)
	ifc.Append(pageStack)

	for _, newsItem := range san.NewsItems {
		if srf := steamNewsFragment(ifc, newsItem); srf != nil {
			pageStack.Append(srf)
		}
	}

	if err := ifc.WriteContent(w); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
	}

	//sb := &strings.Builder{}
	//sanvm := view_models.NewSteamNews(san)
	//if err := tmpl.ExecuteTemplate(sb, "steam-news-content", sanvm); err != nil {
	//	http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
	//	return
	//}
	//gaugin_middleware.DefaultHeaders(w)
	//
	//if err := app.RenderSection(id, stencil_app.SteamNewsSection, sb.String(), w); err != nil {
	//	http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
	//	return
	//}
}

func steamNewsFragment(r compton.Registrar, item steam_integration.NewsItem) compton.Element {

	dsTitle := fspan.Text(r, item.Title).FontWeight(weight.Bolder)
	ds := details_summary.Smaller(r, dsTitle, false)

	return ds
}

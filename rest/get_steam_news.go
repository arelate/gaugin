package rest

import (
	"github.com/arelate/gaugin/rest/compton_data"
	"github.com/arelate/gaugin/rest/compton_fragments"
	"github.com/arelate/gaugin/rest/gaugin_styles"
	"github.com/arelate/southern_light/steam_integration"
	"github.com/boggydigital/compton/consts/color"
	"github.com/boggydigital/compton/consts/direction"
	"github.com/boggydigital/compton/elements/els"
	"github.com/boggydigital/compton/elements/flex_items"
	"github.com/boggydigital/compton/elements/fspan"
	"github.com/boggydigital/compton/elements/iframe_expand"
	"github.com/boggydigital/nod"
	"net/http"
)

func GetSteamNews(w http.ResponseWriter, r *http.Request) {

	// GET /steam-news?id

	id := r.URL.Query().Get("id")
	all := r.URL.Query().Has("all")

	san, err := getSteamNews(http.DefaultClient, id)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	section := compton_data.SteamNewsSection
	ifc := iframe_expand.IframeExpandContent(section, compton_data.SectionTitles[section]).
		AppendStyle(gaugin_styles.SteamNews)

	pageStack := flex_items.FlexItems(ifc, direction.Column)
	ifc.Append(pageStack)

	communityAnnouncements := make([]steam_integration.NewsItem, 0, len(san.NewsItems))
	for _, ni := range san.NewsItems {
		if ni.FeedType != compton_data.FeedTypeCommunityAnnouncement {
			continue
		}
		communityAnnouncements = append(communityAnnouncements, ni)
	}

	if len(san.NewsItems) > 0 &&
		len(communityAnnouncements) < len(san.NewsItems) {
		title := "Show all news items types"
		href := "/steam-news?id=" + id + "&all"
		if all {
			title = "Show only community announcements"
			href = "/steam-news?id=" + id
		}
		pageStack.Append(compton_fragments.ShowMoreButton(ifc, title, href))
	}

	newsItems := communityAnnouncements
	if all {
		newsItems = san.NewsItems
	}

	if len(newsItems) == 0 {
		title := "Community announcements are not available for this product"
		if all {
			title = "Steam news are not available for this product"
		}
		fs := fspan.Text(ifc, title).ForegroundColor(color.Gray)
		pageStack.Append(flex_items.Center(ifc, fs))
	}

	for ii, ni := range newsItems {
		if srf := compton_fragments.SteamNewsItem(ifc, ni, ii == 0); srf != nil {
			pageStack.Append(srf)
		}
		if ii < len(newsItems)-1 {
			pageStack.Append(els.Hr())
		}
	}

	if err := ifc.WriteContent(w); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
	}
}

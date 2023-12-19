package rest

import (
	"fmt"
	"github.com/arelate/gaugin/data"
	"github.com/arelate/gaugin/gaugin_middleware"
	"github.com/arelate/gaugin/paths"
	"github.com/arelate/gaugin/stencil_app"
	"github.com/arelate/southern_light"
	"github.com/arelate/southern_light/gog_integration"
	"github.com/arelate/southern_light/gogdb_integration"
	"github.com/arelate/southern_light/hltb_integration"
	"github.com/arelate/southern_light/igdb_integration"
	"github.com/arelate/southern_light/ign_integration"
	"github.com/arelate/southern_light/mobygames_integration"
	"github.com/arelate/southern_light/pcgw_integration"
	"github.com/arelate/southern_light/protondb_integration"
	"github.com/arelate/southern_light/steam_integration"
	"github.com/arelate/southern_light/strategywiki_integration"
	"github.com/arelate/southern_light/vndb_integration"
	"github.com/arelate/southern_light/wikipedia_integration"
	"github.com/arelate/southern_light/winehq_integration"
	"golang.org/x/exp/maps"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"

	"github.com/arelate/gaugin/view_models"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
)

func GetProduct(w http.ResponseWriter, r *http.Request) {

	// GET /product?slug -> /product?id

	if r.URL.Query().Has(vangogh_local_data.SlugProperty) {
		if ids, _, err := getSearch(http.DefaultClient, r.URL.Query()); err == nil {
			if len(ids) > 0 {
				for _, id := range ids {
					http.Redirect(w, r, paths.ProductId(id), http.StatusPermanentRedirect)
					return
				}
			} else {
				http.Error(w, nod.ErrorStr("unknown slug"), http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
			return
		}
	}

	id := r.URL.Query().Get(vangogh_local_data.IdProperty)

	st := gaugin_middleware.NewServerTimings()
	start := time.Now()

	idRedux, cached, err := getRedux(http.DefaultClient, id, false, stencil_app.ProductProperties...)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	if cached {
		st.SetFlag("getRedux-cached")
	}
	st.Set("getRedux", time.Since(start).Milliseconds())

	// fill redux, data presence to allow showing only the section that will have data

	start = time.Now()
	hasRedux, cached, err := getHasRedux(http.DefaultClient,
		id,
		vangogh_local_data.DescriptionOverviewProperty,
		vangogh_local_data.ChangelogProperty,
		vangogh_local_data.ScreenshotsProperty,
		vangogh_local_data.VideoIdProperty)

	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	if cached {
		st.SetFlag("getHasRedux-cached")
	}
	st.Set("getHasRedux", time.Since(start).Milliseconds())

	hasSections := make([]string, 0)

	if rdx, ok := hasRedux[id]; ok {
		if view_models.FlagFromRedux(rdx, vangogh_local_data.DescriptionOverviewProperty) {
			hasSections = append(hasSections, stencil_app.DescriptionSection)
		}
		if view_models.FlagFromRedux(rdx, vangogh_local_data.ScreenshotsProperty) {
			hasSections = append(hasSections, stencil_app.ScreenshotsSection)
		}
		if view_models.FlagFromRedux(rdx, vangogh_local_data.VideoIdProperty) {
			hasSections = append(hasSections, stencil_app.VideosSection)
		}
		if view_models.FlagFromRedux(rdx, vangogh_local_data.ChangelogProperty) {
			hasSections = append(hasSections, stencil_app.ChangelogSection)
		}
	}

	start = time.Now()
	hasData, cached, err := getHasData(
		http.DefaultClient,
		id,
		vangogh_local_data.SteamAppNews,
		vangogh_local_data.SteamReviews,
		vangogh_local_data.Details)

	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	if cached {
		st.SetFlag("getHasData-cached")
	}
	st.Set("getHasData", time.Since(start).Milliseconds())

	if hasData[vangogh_local_data.SteamAppNews.String()][id] == vangogh_local_data.TrueValue {
		hasSections = append(hasSections, stencil_app.SteamNewsSection)
	}
	if hasData[vangogh_local_data.SteamReviews.String()][id] == vangogh_local_data.TrueValue {
		hasSections = append(hasSections, stencil_app.SteamReviewsSection)
	}
	if hasData[vangogh_local_data.Details.String()][id] == vangogh_local_data.TrueValue {
		hasSections = append(hasSections, stencil_app.DownloadsSection)
	}

	insertAggregateLinks(idRedux[id], id)

	irap := NewIRAProxy(idRedux)

	gaugin_middleware.DefaultHeaders(st, w)

	// adding titles for related games
	relatedIds := make(map[string]bool)
	relatedProps := []string{vangogh_local_data.RequiresGamesProperty, vangogh_local_data.IsRequiredByGamesProperty, vangogh_local_data.IncludesGamesProperty, vangogh_local_data.IsIncludedByGamesProperty}
	for _, p := range relatedProps {
		if rids, ok := irap.GetAllValues(p, id); ok {
			for _, rid := range rids {
				relatedIds[rid] = true
			}
		}
	}
	rids := maps.Keys(relatedIds)
	sort.Strings(rids)
	titleRedux, cached, err := getTitles(http.DefaultClient, rids...)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}
	if cached {
		st.SetFlag("getTitles-cached")
	}
	irap.Merge(titleRedux)

	// adding tag names for related games
	tagNamesRedux, cached, err := getRedux(http.DefaultClient, "", true, vangogh_local_data.TagNameProperty)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}
	if cached {
		st.SetFlag("getRedux-tagNames-cached")
	}
	irap.Merge(tagNamesRedux)

	if err := app.RenderItem(id, hasSections, irap, w); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}
}

func gogLink(p string) string {
	u := url.URL{
		Scheme: southern_light.HttpsScheme,
		Host:   gog_integration.WwwGogHost,
		Path:   p,
	}
	return u.String()
}

func insertAggregateLinks(rdx map[string][]string, id string) {
	for _, p := range []string{
		vangogh_local_data.StoreUrlProperty,
		vangogh_local_data.ForumUrlProperty,
		vangogh_local_data.SupportUrlProperty} {
		if len(rdx[p]) > 0 {
			rdx[data.GauginGOGLinksProperty] = append(rdx[data.GauginGOGLinksProperty],
				//fmt.Sprintf("%s (%s)", p, gogLink(rdx[p][0])))
				fmt.Sprintf("%s=%s", p, gogLink(rdx[p][0])))
		}
	}

	if len(rdx[vangogh_local_data.SteamAppIdProperty]) > 0 {
		if steamAppId := rdx[vangogh_local_data.SteamAppIdProperty][0]; steamAppId != "" {
			if appId, err := strconv.ParseUint(steamAppId, 10, 32); err == nil {
				uAppId := uint32(appId)
				rdx[data.GauginSteamLinksProperty] =
					append(rdx[data.GauginSteamLinksProperty],
						//fmt.Sprintf("%s (%s)", data.GauginSteamCommunityUrlProperty, steam_integration.SteamCommunityUrl(uAppId)))
						fmt.Sprintf("%s=%s", data.GauginSteamCommunityUrlProperty, steam_integration.SteamCommunityUrl(uAppId)))
				rdx[data.GauginOtherLinksProperty] =
					append(rdx[data.GauginOtherLinksProperty],
						//fmt.Sprintf("%s (%s)", data.GauginProtonDBUrlProperty, protondb_integration.ProtonDBUrl(uAppId)))
						fmt.Sprintf("%s=%s", data.GauginProtonDBUrlProperty, protondb_integration.ProtonDBUrl(uAppId)))
			}
		}
	}

	rdx[data.GauginOtherLinksProperty] = append(rdx[data.GauginOtherLinksProperty],
		//fmt.Sprintf("%s (%s)", data.GauginGOGDBUrlProperty, gogdb_integration.GOGDBUrl(id)))
		fmt.Sprintf("%s=%s", data.GauginGOGDBUrlProperty, gogdb_integration.GOGDBUrl(id)))

	otherLink(rdx,
		vangogh_local_data.PCGWPageIdProperty,
		data.GauginPCGamingWikiUrlProperty,
		pcgw_integration.WikiUrl)
	otherLink(rdx,
		vangogh_local_data.HLTBIdProperty,
		data.GauginHLTBUrlProperty,
		hltb_integration.GameUrl)
	otherLink(rdx,
		vangogh_local_data.IGDBIdProperty,
		data.GauginIGDBUrlProperty,
		igdb_integration.GameUrl)
	otherLink(rdx,
		vangogh_local_data.StrategyWikiIdProperty,
		data.GauginStrategyWikiUrlProperty,
		strategywiki_integration.WikiUrl)
	otherLink(rdx,
		vangogh_local_data.MobyGamesIdProperty,
		data.GauginMobyGamesUrlProperty,
		mobygames_integration.GameUrl)
	otherLink(rdx,
		vangogh_local_data.WikipediaIdProperty,
		data.GauginWikipediaUrlProperty,
		wikipedia_integration.WikiUrl)
	otherLink(rdx,
		vangogh_local_data.WineHQIdProperty,
		data.GauginWineHQUrlProperty,
		winehq_integration.WineHQUrl)
	otherLink(rdx,
		vangogh_local_data.VNDBIdProperty,
		data.GauginVNDBUrlProperty,
		vndb_integration.ItemUrl)
	otherLink(rdx,
		vangogh_local_data.IGNWikiSlugProperty,
		data.GauginIGNWikiUrlProperty,
		ign_integration.WikiUrl)
}

func otherLink(rdx map[string][]string, p string, up string, uf func(string) *url.URL) {
	if len(rdx[p]) > 0 {
		id := rdx[p][0]
		rdx[data.GauginOtherLinksProperty] = append(rdx[data.GauginOtherLinksProperty],
			//fmt.Sprintf("%s (%s)", up, uf(id)))
			fmt.Sprintf("%s=%s", up, uf(id)))
	}
}

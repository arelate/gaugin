package rest

import (
	"fmt"
	"github.com/arelate/gaugin/data"
	"github.com/arelate/gaugin/paths"
	"github.com/arelate/gaugin/rest/compton_data"
	"github.com/arelate/gaugin/rest/compton_pages"
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
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/kevlar"
	"github.com/boggydigital/nod"
	"golang.org/x/exp/maps"
	"net/http"
	"net/url"
	"sort"
	"strconv"
)

var (
	propertiesSections = map[string]string{
		vangogh_local_data.DescriptionOverviewProperty: compton_data.DescriptionSection,
		vangogh_local_data.ChangelogProperty:           compton_data.ChangelogSection,
		vangogh_local_data.ScreenshotsProperty:         compton_data.ScreenshotsSection,
		vangogh_local_data.VideoIdProperty:             compton_data.VideosSection,
	}
	propertiesSectionsOrder = []string{
		vangogh_local_data.DescriptionOverviewProperty,
		vangogh_local_data.ChangelogProperty,
		vangogh_local_data.ScreenshotsProperty,
		vangogh_local_data.VideoIdProperty,
	}

	dataTypesSections = map[vangogh_local_data.ProductType]string{
		vangogh_local_data.SteamAppNews:                 compton_data.SteamNewsSection,
		vangogh_local_data.SteamReviews:                 compton_data.SteamReviewsSection,
		vangogh_local_data.SteamDeckCompatibilityReport: compton_data.SteamDeckSection,
		vangogh_local_data.Details:                      compton_data.DownloadsSection,
	}

	dataTypesSectionsOrder = []vangogh_local_data.ProductType{
		vangogh_local_data.SteamAppNews,
		vangogh_local_data.SteamReviews,
		vangogh_local_data.SteamDeckCompatibilityReport,
		vangogh_local_data.Details,
	}
)

func propertiesFromRedux(redux map[string][]string, property string) []string {
	if val, ok := redux[property]; ok {
		return val
	} else {
		return []string{}
	}
}

func propertyFromRedux(redux map[string][]string, property string) string {
	properties := propertiesFromRedux(redux, property)
	if len(properties) > 0 {
		return properties[0]
	}
	return ""
}

// TODO: review and remove this
func FlagFromRedux(redux map[string][]string, property string) bool {
	return propertyFromRedux(redux, property) == vangogh_local_data.TrueValue
}

func GetProduct(w http.ResponseWriter, r *http.Request) {

	// GET /product?slug -> /product?id

	if r.URL.Query().Has(vangogh_local_data.SlugProperty) {
		if ids, err := getSearch(http.DefaultClient, r.URL.Query()); err == nil {
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

	idRedux, err := getRedux(http.DefaultClient, id, false, compton_data.ProductProperties...)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	// fill redux, data presence to allow showing only the section that will have data

	hasRedux, err := getHasRedux(http.DefaultClient, id, maps.Keys(propertiesSections)...)

	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	hasSections := make([]string, 0)
	// every product is expected to have at least those sections
	hasSections = append(hasSections, compton_data.PropertiesSection, compton_data.ExternalLinksSection)

	if hRdx, ok := hasRedux[id]; ok {
		for _, property := range propertiesSectionsOrder {
			if section, ok := propertiesSections[property]; ok {
				if FlagFromRedux(hRdx, property) {
					hasSections = append(hasSections, section)
				}
			}
		}
	}

	hasData, err := getHasData(http.DefaultClient, id, maps.Keys(dataTypesSections)...)

	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	for _, dt := range dataTypesSectionsOrder {
		if section, ok := dataTypesSections[dt]; ok {
			if hasData[dt.String()][id] == vangogh_local_data.TrueValue {
				hasSections = append(hasSections, section)
			}
		}
	}

	insertAggregateLinks(idRedux[id], id)

	//gaugin_middleware.DefaultHeaders(w)

	// adding titles for related games
	relatedIds := make(map[string]bool)
	relatedProps := []string{
		vangogh_local_data.RequiresGamesProperty,
		vangogh_local_data.IsRequiredByGamesProperty,
		vangogh_local_data.IncludesGamesProperty,
		vangogh_local_data.IsIncludedByGamesProperty}
	for _, p := range relatedProps {
		if pvs, ok := idRedux[id]; ok {
			if rids, ok := pvs[p]; ok {
				for _, rid := range rids {
					relatedIds[rid] = true
				}
			}
		}
	}
	rids := maps.Keys(relatedIds)
	sort.Strings(rids)
	titleRedux, err := getTitles(http.DefaultClient, rids...)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	idRedux = MergeIdPropertyValues(idRedux, titleRedux)

	// adding tag names for related games
	tagNamesRedux, err := getRedux(http.DefaultClient, "", true, vangogh_local_data.TagNameProperty)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}
	rdx := kevlar.ReduxProxy(MergeIdPropertyValues(idRedux, tagNamesRedux))

	if productPage := compton_pages.Product(id, rdx, hasSections); productPage != nil {
		if err := productPage.WriteContent(w); err != nil {
			http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		}
	} else {
		http.NotFound(w, r)
	}

	//if err := app.RenderItem(id, hasSections, rdx, w); err != nil {
	//	http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
	//	return
	//}
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
				fmt.Sprintf("%s=%s", p, gogLink(rdx[p][0])))
		}
	}

	if len(rdx[vangogh_local_data.SteamAppIdProperty]) > 0 {
		if steamAppId := rdx[vangogh_local_data.SteamAppIdProperty][0]; steamAppId != "" {
			if appId, err := strconv.ParseUint(steamAppId, 10, 32); err == nil {
				uAppId := uint32(appId)
				rdx[data.GauginSteamLinksProperty] =
					append(rdx[data.GauginSteamLinksProperty],
						fmt.Sprintf("%s=%s", data.GauginSteamCommunityUrlProperty, steam_integration.SteamCommunityUrl(uAppId)))
				rdx[data.GauginOtherLinksProperty] =
					append(rdx[data.GauginOtherLinksProperty],
						fmt.Sprintf("%s=%s", data.GauginProtonDBUrlProperty, protondb_integration.ProtonDBUrl(uAppId)))
			}
		}
	}

	rdx[data.GauginOtherLinksProperty] = append(rdx[data.GauginOtherLinksProperty],
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

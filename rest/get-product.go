package rest

import (
	"github.com/arelate/gaugin/gaugin_middleware"
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/http"
)

var productProperties = []string{
	vangogh_local_data.DehydratedImageProperty,
	vangogh_local_data.ImageProperty,
	vangogh_local_data.ProductTypeProperty,
	vangogh_local_data.TitleProperty,
	vangogh_local_data.TagIdProperty,
	vangogh_local_data.LocalTagsProperty,
	vangogh_local_data.OperatingSystemsProperty,
	vangogh_local_data.RatingProperty,
	vangogh_local_data.DevelopersProperty,
	vangogh_local_data.PublisherProperty,
	vangogh_local_data.SeriesProperty,
	vangogh_local_data.GenresProperty,
	vangogh_local_data.PropertiesProperty,
	vangogh_local_data.FeaturesProperty,
	vangogh_local_data.LanguageCodeProperty,
	vangogh_local_data.GlobalReleaseDateProperty,
	vangogh_local_data.GOGReleaseDateProperty,
	vangogh_local_data.GOGOrderDateProperty,
	vangogh_local_data.IncludesGamesProperty,
	vangogh_local_data.IsIncludedByGamesProperty,
	vangogh_local_data.RequiresGamesProperty,
	vangogh_local_data.IsRequiredByGamesProperty,
	vangogh_local_data.StoreUrlProperty,
	vangogh_local_data.ForumUrlProperty,
	vangogh_local_data.SupportUrlProperty,
	vangogh_local_data.WishlistedProperty,
	vangogh_local_data.OwnedProperty,
	vangogh_local_data.IsFreeProperty,
	vangogh_local_data.IsDiscountedProperty,
	vangogh_local_data.PreOrderProperty,
	vangogh_local_data.TBAProperty,
	vangogh_local_data.ComingSoonProperty,
	vangogh_local_data.InDevelopmentProperty,
	vangogh_local_data.IsUsingDOSBoxProperty,
	vangogh_local_data.IsUsingScummVMProperty,
	vangogh_local_data.BasePriceProperty,
	vangogh_local_data.PriceProperty,
	vangogh_local_data.DiscountPercentageProperty,
	vangogh_local_data.SteamAppIdProperty,
	vangogh_local_data.SteamReviewScoreDescProperty,
	vangogh_local_data.SteamTagsProperty,
}

//National flags that correspond to language code.
//In some cases there is no obvious way to map those,
//attempting to use sensible option: ar, ca, fa
//Few options are not possible to map to countries (left as comments below)
var languageFlags = map[string]string{
	"en":    "🇺🇸",
	"de":    "🇩🇪",
	"fr":    "🇫🇷",
	"es":    "🇪🇸",
	"ru":    "🇷🇺",
	"it":    "🇮🇹",
	"cn":    "🇨🇳",
	"jp":    "🇯🇵",
	"pl":    "🇵🇱",
	"br":    "🇧🇷",
	"ko":    "🇰🇷",
	"zh":    "🇨🇳",
	"tr":    "🇹🇷",
	"cz":    "🇨🇿",
	"pt":    "🇵🇹",
	"nl":    "🇳🇱",
	"es_mx": "🇲🇽",
	"hu":    "🇭🇺",
	"uk":    "🇺🇦",
	"ar":    "🇸🇦",
	"sv":    "🇸🇪",
	"no":    "🇳🇴",
	"da":    "🇩🇰",
	"fi":    "🇫🇮",
	"th":    "🇹🇭",
	"ro":    "🇷🇴",
	"gk":    "🇬🇷",
	"bl":    "🇧🇬",
	"sk":    "🇸🇮",
	"be":    "🇧🇾",
	"he":    "🇮🇱",
	"sb":    "🇷🇸",
	"ca":    "🇪🇸",
	"is":    "🇮🇸",
	"fa":    "🇮🇷",
	"et":    "🇪🇪",
	//Inuktitut (gog_IN): 2 items
	//latine (la): 1 items
}

func GetProduct(w http.ResponseWriter, r *http.Request) {

	// GET /product?slug -> /product?id

	if r.URL.Query().Has(vangogh_local_data.SlugProperty) {
		if idSet, err := vangogh_local_data.IdSetFromUrl(r.URL); err == nil {
			if len(idSet) > 0 {
				for id := range idSet {
					http.Redirect(w, r, "/product?id="+id, http.StatusPermanentRedirect)
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

	idRedux, err := getRedux(http.DefaultClient, id, productProperties...)
	if err != nil {
		http.Error(w, nod.ErrorStr("error getting redux"), http.StatusInternalServerError)
		return
	}

	gaugin_middleware.DefaultHeaders(w)

	pvm, err := productViewModelFromRedux(idRedux)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	// fill redux, data presence to allow showing only the section that will have data

	hasData, err := getHasData(
		http.DefaultClient,
		id,
		gog_integration.Game,
		vangogh_local_data.SteamAppNews,
		vangogh_local_data.Details)

	if err != nil {
		http.Error(w, nod.ErrorStr("error getting has_data"), http.StatusInternalServerError)
		return
	}

	pvm.HasSteamAppNews = hasData[vangogh_local_data.SteamAppNews.String()][id] == vangogh_local_data.TrueValue
	pvm.HasDownloads = hasData[vangogh_local_data.Details.String()][id] == vangogh_local_data.TrueValue

	hasRedux, err := getHasRedux(http.DefaultClient, id,
		vangogh_local_data.DescriptionOverviewProperty,
		vangogh_local_data.ChangelogProperty,
		vangogh_local_data.ScreenshotsProperty,
		vangogh_local_data.VideoIdProperty)

	if err != nil {
		http.Error(w, nod.ErrorStr("error getting has_redux"), http.StatusInternalServerError)
		return
	}

	if rdx, ok := hasRedux[id]; ok {
		pvm.HasDescription = flagFromRedux(rdx, vangogh_local_data.DescriptionOverviewProperty)
		pvm.HasChangelog = flagFromRedux(rdx, vangogh_local_data.ChangelogProperty)
		pvm.HasScreenshots = flagFromRedux(rdx, vangogh_local_data.ScreenshotsProperty)
		pvm.HasVideos = flagFromRedux(rdx, vangogh_local_data.VideoIdProperty)
	}

	if err := tmpl.ExecuteTemplate(w, "product-page", pvm); err != nil {
		http.Error(w, nod.ErrorStr("template exec error"), http.StatusInternalServerError)
		return
	}
}

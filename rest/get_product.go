package rest

import (
	"github.com/arelate/gaugin/gaugin_middleware"
	"github.com/arelate/gaugin/stencil_app"
	"net/http"
	"time"

	"github.com/arelate/gaugin/view_models"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
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
	vangogh_local_data.PublishersProperty,
	vangogh_local_data.SeriesProperty,
	vangogh_local_data.GenresProperty,
	vangogh_local_data.StoreTagsProperty,
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
	vangogh_local_data.ValidationResultProperty,
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

	st := gaugin_middleware.NewServerTimings()
	start := time.Now()

	idRedux, cached, err := getRedux(http.DefaultClient, id, false, productProperties...)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	if cached {
		st.SetFlag("getRedux-cached")
	}
	st.Set("getRedux", time.Since(start).Milliseconds())

	pvm, err := view_models.NewProduct(idRedux)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

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

	if rdx, ok := hasRedux[id]; ok {
		if view_models.FlagFromRedux(rdx, vangogh_local_data.DescriptionOverviewProperty) {
			pvm.Sections = append(pvm.Sections, stencil_app.DescriptionSection)
		}
		if view_models.FlagFromRedux(rdx, vangogh_local_data.ScreenshotsProperty) {
			pvm.Sections = append(pvm.Sections, stencil_app.ScreenshotsSection)
		}
		if view_models.FlagFromRedux(rdx, vangogh_local_data.VideoIdProperty) {
			pvm.Sections = append(pvm.Sections, stencil_app.VideosSection)
		}
		if view_models.FlagFromRedux(rdx, vangogh_local_data.ChangelogProperty) {
			pvm.Sections = append(pvm.Sections, stencil_app.ChangelogSection)
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
		pvm.Sections = append(pvm.Sections, stencil_app.SteamNewsSection)
	}
	if hasData[vangogh_local_data.SteamReviews.String()][id] == vangogh_local_data.TrueValue {
		pvm.Sections = append(pvm.Sections, stencil_app.SteamReviewsSection)
	}
	if hasData[vangogh_local_data.Details.String()][id] == vangogh_local_data.TrueValue {
		pvm.Sections = append(pvm.Sections, stencil_app.DownloadsSection)
	}

	gaugin_middleware.DefaultHeaders(st, w)

	if err := tmpl.ExecuteTemplate(w, "product-page", pvm); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}
}

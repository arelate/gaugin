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

	idRedux, cached, err := getRedux(http.DefaultClient, id, false, stencil_app.ProductProperties...)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	if cached {
		st.SetFlag("getRedux-cached")
	}
	st.Set("getRedux", time.Since(start).Milliseconds())

	//pvm, err := view_models.NewProduct(idRedux)
	//if err != nil {
	//	http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
	//	return
	//}

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

	gaugin_middleware.DefaultHeaders(st, w)

	irap := vangogh_local_data.NewIRAProxy(idRedux)

	if err := app.RenderItem(id, hasSections, irap, w); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	//if err := tmpl.ExecuteTemplate(w, "product-page", pvm); err != nil {
	//	http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
	//	return
	//}
}

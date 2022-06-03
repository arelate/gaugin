package rest

import (
	"github.com/arelate/gaugin/gaugin_middleware"
	"net/http"
	"strings"

	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
)

var productProperties = []string{
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
	vangogh_local_data.ScreenshotsProperty,
	vangogh_local_data.VideoIdProperty,
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
}

func GetProduct(w http.ResponseWriter, r *http.Request) {

	// GET /product?slug -> /product?id

	if r.URL.Query().Has("slug") {
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

	id := r.URL.Query().Get("id")

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

	validVideos := make([]string, 0, len(pvm.Videos))
	for _, v := range pvm.Videos {
		if strings.Contains(v, "(") {
			continue
		}
		validVideos = append(validVideos, v)
	}
	pvm.Videos = validVideos

	if err := getCurrentOtherOSDownloads(pvm, id, r.Header.Get("User-Agent")); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "product-page", pvm); err != nil {
		http.Error(w, nod.ErrorStr("template exec error"), http.StatusInternalServerError)
		return
	}
}

func getCurrentOtherOSDownloads(pvm *productViewModel, id string, userAgent string) error {

	if !pvm.Owned {
		return nil
	}

	dls, err := getDownloads(http.DefaultClient, id, operatingSystems, languageCodes)
	if err != nil {
		return err
	}

	var currentOS vangogh_local_data.OperatingSystem
	if strings.Contains(userAgent, "Windows") {
		currentOS = vangogh_local_data.Windows
	} else if strings.Contains(userAgent, "Mac OS X") {
		currentOS = vangogh_local_data.MacOS
	} else if strings.Contains(userAgent, "Linux") {
		currentOS = vangogh_local_data.Linux
	}

	pvm.CurrentOS = &productDownloads{
		CurrentOS:  true,
		Installers: make(vangogh_local_data.DownloadsList, 0, len(dls)),
		DLCs:       make(vangogh_local_data.DownloadsList, 0, len(dls)),
		Extras:     make(vangogh_local_data.DownloadsList, 0, len(dls)),
	}
	pvm.OtherOS = &productDownloads{
		CurrentOS:  false,
		Installers: make(vangogh_local_data.DownloadsList, 0, len(dls)),
		DLCs:       make(vangogh_local_data.DownloadsList, 0, len(dls)),
		Extras:     make(vangogh_local_data.DownloadsList, 0, len(dls)),
	}

	var osd *productDownloads
	for _, dl := range dls {
		if dl.OS == currentOS ||
			dl.OS == vangogh_local_data.AnyOperatingSystem {
			osd = pvm.CurrentOS
		} else {
			osd = pvm.OtherOS
		}

		switch dl.Type {
		case vangogh_local_data.Installer:
			osd.Installers = append(osd.Installers, dl)
		case vangogh_local_data.DLC:
			osd.DLCs = append(osd.DLCs, dl)
		case vangogh_local_data.Extra:
			fallthrough
		default:
			osd.Extras = append(osd.Extras, dl)
		}
	}

	return nil
}

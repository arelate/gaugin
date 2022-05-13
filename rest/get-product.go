package rest

import (
	"net/http"
	"strings"

	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
)

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

	idRedux, err := getRedux(http.DefaultClient, id, vangogh_local_data.ReduxProperties()...)
	if err != nil {
		http.Error(w, nod.ErrorStr("error getting redux"), http.StatusInternalServerError)
		return
	}

	defaultHeaders(w)

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

	if pvm.Owned {
		dls, err := getDownloads(http.DefaultClient, id, operatingSystems, languageCodes)
		if err != nil {
			http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
			return
		}

		var currentOS vangogh_local_data.OperatingSystem
		userAgent := r.Header.Get("User-Agent")
		if strings.Contains(userAgent, "Windows") {
			currentOS = vangogh_local_data.Windows
		} else if strings.Contains(userAgent, "Mac OS X") {
			currentOS = vangogh_local_data.MacOS
		} else if strings.Contains(userAgent, "Linux") {
			currentOS = vangogh_local_data.Linux
		}

		pvm.CurrentOS = &productDownloads{
			Context:    "current-operating-system",
			Installers: make(vangogh_local_data.DownloadsList, 0, len(dls)),
			DLCs:       make(vangogh_local_data.DownloadsList, 0, len(dls)),
			Extras:     make(vangogh_local_data.DownloadsList, 0, len(dls)),
		}
		pvm.OtherOS = &productDownloads{
			Context:    "other-operating-system",
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
	}

	if err := tmpl.ExecuteTemplate(w, "product", pvm); err != nil {
		http.Error(w, nod.ErrorStr("template exec error"), http.StatusInternalServerError)
		return
	}
}

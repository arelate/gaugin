package rest

import (
	"github.com/arelate/gaugin/gaugin_middleware"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"golang.org/x/exp/maps"
	"net/http"
	"strings"
)

func GetDownloads(w http.ResponseWriter, r *http.Request) {

	// GET /downloads?id

	id := r.URL.Query().Get("id")

	gaugin_middleware.DefaultHeaders(w)

	dvm, err := getCurrentOtherOSDownloads(id, r.Header.Get("User-Agent"))
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "downloads-page", dvm); err != nil {
		http.Error(w, nod.ErrorStr("template exec error"), http.StatusInternalServerError)
		return
	}
}

func getCurrentOtherOSDownloads(id string, userAgent string) (*downloadsViewModel, error) {

	dvm := &downloadsViewModel{
		Context: "iframe",
	}

	//details, err := getDetails(http.DefaultClient, id)
	//fmt.Println(details.Title)

	dls, err := getDownloads(http.DefaultClient, id, operatingSystems, languageCodes)
	if err != nil {
		return dvm, err
	}

	var currentOS vangogh_local_data.OperatingSystem
	if strings.Contains(userAgent, "Windows") {
		currentOS = vangogh_local_data.Windows
	} else if strings.Contains(userAgent, "Mac OS X") {
		currentOS = vangogh_local_data.MacOS
	} else if strings.Contains(userAgent, "Linux") {
		currentOS = vangogh_local_data.Linux
	}

	otherOS := make(map[string]interface{}, 0)

	dvm.CurrentOS = &productDownloads{
		OperatingSystems: currentOS.String(),
		CurrentOS:        true,
		Installers:       make(vangogh_local_data.DownloadsList, 0, len(dls)),
		DLCs:             make(vangogh_local_data.DownloadsList, 0, len(dls)),
		Extras:           make(vangogh_local_data.DownloadsList, 0, len(dls)),
	}
	dvm.OtherOS = &productDownloads{
		CurrentOS:  false,
		Installers: make(vangogh_local_data.DownloadsList, 0, len(dls)),
		DLCs:       make(vangogh_local_data.DownloadsList, 0, len(dls)),
		Extras:     make(vangogh_local_data.DownloadsList, 0, len(dls)),
	}

	var osd *productDownloads
	for _, dl := range dls {
		if dl.OS == currentOS ||
			dl.OS == vangogh_local_data.AnyOperatingSystem {
			osd = dvm.CurrentOS
		} else {
			otherOS[dl.OS.String()] = nil
			osd = dvm.OtherOS
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

	dvm.OtherOS.OperatingSystems = strings.Join(maps.Keys(otherOS), ", ")

	return dvm, nil
}

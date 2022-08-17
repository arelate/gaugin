package rest

import (
	"net/http"
	"strings"

	"github.com/arelate/gaugin/gaugin_middleware"
	"github.com/arelate/gaugin/view_models"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
)

func GetDownloads(w http.ResponseWriter, r *http.Request) {

	// GET /downloads?id

	id := r.URL.Query().Get("id")

	gaugin_middleware.DefaultHeaders(w)

	clientOS := getClientOperatingSystem(r)
	dvm, err := getDownloadsViewModel(id, clientOS)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "downloads-page", dvm); err != nil {
		http.Error(w, nod.ErrorStr("template exec error"), http.StatusInternalServerError)
		return
	}
}

func getClientOperatingSystem(r *http.Request) vangogh_local_data.OperatingSystem {

	var clientOS vangogh_local_data.OperatingSystem

	//attempt to extract platform from user agent client hints first
	secChUaPlatform := r.Header.Get("Sec-CH-UA-Platform")

	switch secChUaPlatform {
	case "Linux":
		clientOS = vangogh_local_data.Linux
	case "iOS":
		fallthrough
	case "macOS":
		clientOS = vangogh_local_data.MacOS
	case "Windows":
		clientOS = vangogh_local_data.Windows
	default:
		// "Android", "Chrome OS", "Chromium OS" or "Unknown"
		clientOS = vangogh_local_data.AnyOperatingSystem
	}

	if clientOS != vangogh_local_data.AnyOperatingSystem {
		return clientOS
	}

	//use "User-Agent" header if we couldn't extract platform from user agent client hints
	userAgent := r.UserAgent()

	if strings.Contains(userAgent, "Windows") {
		clientOS = vangogh_local_data.Windows
	} else if strings.Contains(userAgent, "Mac OS X") {
		clientOS = vangogh_local_data.MacOS
	} else if strings.Contains(userAgent, "Linux") {
		clientOS = vangogh_local_data.Linux
	}

	return clientOS
}

func getDownloadsViewModel(id string, clientOS vangogh_local_data.OperatingSystem) (*view_models.Downloads, error) {

	//we specifically get /downloads and not /data&product-type=details because of Details
	//format complexities, see gog_integration/details.go/GetGameDownloads comment
	dls, err := getDownloads(http.DefaultClient, id, operatingSystems, languageCodes)
	if err != nil {
		return nil, err
	}

	dvm := view_models.NewDownloads(clientOS, dls)

	return dvm, nil
}

package rest

import (
	"github.com/arelate/gaugin/rest/compton_pages"
	"github.com/boggydigital/kevlar"
	"net/http"
	"strings"

	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
)

func GetDownloads(w http.ResponseWriter, r *http.Request) {

	// GET /downloads?id

	id := r.URL.Query().Get("id")

	idRedux, err := getRedux(
		http.DefaultClient,
		id,
		false,
		vangogh_local_data.ValidationResultProperty,
		vangogh_local_data.ValidationCompletedProperty)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	//we specifically get /downloads and not /data&product-type=details because of Details
	//format complexities, see gog_integration/details.go/GetGameDownloads comment
	dls, err := getDownloads(http.DefaultClient, id, operatingSystems, languageCodes, excludePatches)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	p := compton_pages.Downloads(id, dls, kevlar.ReduxProxy(idRedux))

	if err := p.WriteContent(w); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
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

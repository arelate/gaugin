package rest

import (
	"github.com/arelate/gaugin/gaugin_middleware"
	"github.com/arelate/gaugin/stencil_app"
	"net/http"
	"strings"

	"github.com/arelate/gaugin/view_models"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
)

func GetDownloads(w http.ResponseWriter, r *http.Request) {

	// GET /downloads?id

	id := r.URL.Query().Get("id")

	clientOS := getClientOperatingSystem(r)
	dvm, err := getDownloadsViewModel(id, clientOS)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	sb := &strings.Builder{}

	if err := tmpl.ExecuteTemplate(sb, "downloads-content", dvm); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	gaugin_middleware.DefaultHeaders(nil, w)

	if err := app.RenderSection(id, stencil_app.DownloadsSection, sb.String(), w); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
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

	idRdx, _, err := getRedux(
		http.DefaultClient,
		id,
		false,
		vangogh_local_data.ValidationResultProperty,
		vangogh_local_data.ValidationCompletedProperty)
	if err != nil {
		return nil, err
	}

	//we specifically get /downloads and not /data&product-type=details because of Details
	//format complexities, see gog_integration/details.go/GetGameDownloads comment
	dls, _, err := getDownloads(http.DefaultClient, id, operatingSystems, languageCodes)
	if err != nil {
		return nil, err
	}

	dvm := view_models.NewDownloads(idRdx[id], clientOS, dls)

	return dvm, nil
}

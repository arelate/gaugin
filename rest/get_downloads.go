package rest

import (
	"fmt"
	"github.com/arelate/gaugin/rest/compton_data"
	"github.com/arelate/gaugin/rest/gaugin_styles"
	"github.com/boggydigital/compton/consts/direction"
	"github.com/boggydigital/compton/elements/flex_items"
	"github.com/boggydigital/compton/elements/iframe_expand"
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

	idRdx, err := getRedux(
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

	dvm := view_models.NewDownloads(idRdx[id], clientOS, dls)
	fmt.Println(dvm)

	section := compton_data.DownloadsSection
	ifc := iframe_expand.IframeExpandContent(section, compton_data.SectionTitles[section]).
		AppendStyle(gaugin_styles.DownloadsStyle)

	pageStack := flex_items.FlexItems(ifc, direction.Column)
	ifc.Append(pageStack)

	if err := ifc.WriteContent(w); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
	}

	//sb := &strings.Builder{}
	//
	//if err := tmpl.ExecuteTemplate(sb, "downloads-content", dvm); err != nil {
	//	http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
	//	return
	//}
	//
	//gaugin_middleware.DefaultHeaders(w)
	//
	//if err := app.RenderSection(id, stencil_app.DownloadsSection, sb.String(), w); err != nil {
	//	http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
	//	return
	//}
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

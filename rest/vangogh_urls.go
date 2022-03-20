package rest

import (
	"fmt"
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_local_data"
	"net/url"
	"strconv"
	"strings"
)

var (
	vangoghScheme  = ""
	vangoghAddress = ""
	vangoghPort    = ""
)

func vangoghHost() string {
	return fmt.Sprintf("%s:%s", vangoghAddress, vangoghPort)
}

func SetVangoghConnection(scheme, address string, port int) {
	vangoghScheme = scheme
	vangoghAddress = address
	vangoghPort = strconv.Itoa(port)
}

const (
	cvEndpoint        = "/v1"
	keysEndpoint      = cvEndpoint + "/keys"
	allReduxEndpoint  = cvEndpoint + "/all_redux"
	reduxEndpoint     = cvEndpoint + "/redux"
	searchEndpoint    = cvEndpoint + "/search"
	downloadsEndpoint = cvEndpoint + "/downloads"
	updatesEndpoint   = cvEndpoint + "/updates"
)

func defaultSort(pt vangogh_local_data.ProductType) string {
	switch pt {
	case vangogh_local_data.StoreProducts:
		return vangogh_local_data.GOGReleaseDateProperty
	case vangogh_local_data.Details:
		return vangogh_local_data.GOGOrderDateProperty
	default:
		return vangogh_local_data.TitleProperty
	}
}

func defaultDesc(pt vangogh_local_data.ProductType) string {
	switch pt {
	case vangogh_local_data.StoreProducts:
		return "true"
	case vangogh_local_data.Details:
		return "true"
	default:
		return "false"
	}
}

func keysUrl(pt vangogh_local_data.ProductType, mt gog_integration.Media) *url.URL {
	u := &url.URL{
		Scheme: vangoghScheme,
		Host:   vangoghHost(),
		Path:   keysEndpoint,
	}
	q := u.Query()
	q.Set("product-type", pt.String())
	q.Set("media", mt.String())
	q.Set("sort", defaultSort(pt))
	q.Set("desc", defaultDesc(pt))
	u.RawQuery = q.Encode()

	return u
}

func allReduxUrl(pt vangogh_local_data.ProductType, mt gog_integration.Media, properties ...string) *url.URL {
	u := &url.URL{
		Scheme: vangoghScheme,
		Host:   vangoghHost(),
		Path:   allReduxEndpoint,
	}
	q := u.Query()
	q.Set("product-type", pt.String())
	q.Set("media", mt.String())
	q.Set("property", strings.Join(properties, ","))
	u.RawQuery = q.Encode()

	return u
}

func reduxUrl(id string, properties ...string) *url.URL {
	u := &url.URL{
		Scheme: vangoghScheme,
		Host:   vangoghHost(),
		Path:   reduxEndpoint,
	}
	q := u.Query()
	q.Set("id", id)
	q.Set("property", strings.Join(properties, ","))
	u.RawQuery = q.Encode()

	return u
}

func searchUrl(q url.Values) *url.URL {
	u := &url.URL{
		Scheme: vangoghScheme,
		Host:   vangoghHost(),
		Path:   searchEndpoint,
	}
	u.RawQuery = q.Encode()

	return u
}

func downloadsUrl(
	id string,
	operatingSystems []vangogh_local_data.OperatingSystem,
	languageCodes []string) *url.URL {
	u := &url.URL{
		Scheme: vangoghScheme,
		Host:   vangoghHost(),
		Path:   downloadsEndpoint,
	}
	q := u.Query()
	q.Set("id", id)
	osStr := make([]string, 0, len(operatingSystems))
	for _, os := range operatingSystems {
		osStr = append(osStr, os.String())
	}
	q.Set("operating-system", strings.Join(osStr, ","))
	q.Set("language-code", strings.Join(languageCodes, ","))
	u.RawQuery = q.Encode()

	return u
}

func updatesUrl(
	mt gog_integration.Media,
	since int) *url.URL {
	u := &url.URL{
		Scheme: vangoghScheme,
		Host:   vangoghHost(),
		Path:   updatesEndpoint,
	}
	q := u.Query()
	q.Set("media", mt.String())
	q.Set("since-hours-ago", strconv.Itoa(since))
	u.RawQuery = q.Encode()

	return u
}

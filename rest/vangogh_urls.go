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
	digestEndpoint    = cvEndpoint + "/digest"
	downloadsEndpoint = cvEndpoint + "/downloads"
	reduxEndpoint     = cvEndpoint + "/redux"
	searchEndpoint    = cvEndpoint + "/search"
	updatesEndpoint   = cvEndpoint + "/updates"
	dataEndpoint      = cvEndpoint + "/data"
)

func reduxUrl(id string, properties ...string) *url.URL {
	u := &url.URL{
		Scheme: vangoghScheme,
		Host:   vangoghHost(),
		Path:   reduxEndpoint,
	}
	q := u.Query()
	q.Set(vangogh_local_data.IdProperty, id)
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

func digestUrl(properties ...string) *url.URL {
	u := &url.URL{
		Scheme: vangoghScheme,
		Host:   vangoghHost(),
		Path:   digestEndpoint,
	}
	q := u.Query()
	q.Set("property", strings.Join(properties, ","))
	u.RawQuery = q.Encode()

	return u
}

func steamAppNewsUrl(id string) *url.URL {
	u := &url.URL{
		Scheme: vangoghScheme,
		Host:   vangoghHost(),
		Path:   dataEndpoint,
	}

	q := u.Query()
	q.Set(vangogh_local_data.ProductTypeProperty, vangogh_local_data.SteamAppNews.String())
	q.Set("media", gog_integration.Game.String())
	q.Set(vangogh_local_data.IdProperty, id)
	q.Set("format", "json")
	u.RawQuery = q.Encode()

	return u
}

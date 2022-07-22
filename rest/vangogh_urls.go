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
	dataEndpoint      = "/data"
	digestEndpoint    = "/digest"
	downloadsEndpoint = "/downloads"
	hasDataEndpoint   = "/has_data"
	hasReduxEndpoint  = "/has_redux"
	localTagEndpoint  = "/local_tag"
	reduxEndpoint     = "/redux"
	searchEndpoint    = "/search"
	tagEndpoint       = "/tag"
	wishlistEndpoint  = "/wishlist"
)

func reduxUrl(id string, all bool, properties ...string) *url.URL {
	u := &url.URL{
		Scheme: vangoghScheme,
		Host:   vangoghHost(),
		Path:   reduxEndpoint,
	}
	q := u.Query()
	if id != "" {
		q.Set(vangogh_local_data.IdProperty, id)
	}
	q.Set("property", strings.Join(properties, ","))
	if all {
		q.Set("all", vangogh_local_data.TrueValue)
	}
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

func dataUrl(id string,
	pt vangogh_local_data.ProductType,
	mt gog_integration.Media) *url.URL {
	u := &url.URL{
		Scheme: vangoghScheme,
		Host:   vangoghHost(),
		Path:   dataEndpoint,
	}

	q := u.Query()
	q.Set(vangogh_local_data.ProductTypeProperty, pt.String())
	q.Set("media", mt.String())
	q.Set(vangogh_local_data.IdProperty, id)
	u.RawQuery = q.Encode()

	return u
}

func hasReduxUrl(id string, properties ...string) *url.URL {
	u := &url.URL{
		Scheme: vangoghScheme,
		Host:   vangoghHost(),
		Path:   hasReduxEndpoint,
	}

	q := u.Query()
	q.Set(vangogh_local_data.IdProperty, id)
	q.Set("property", strings.Join(properties, ","))
	u.RawQuery = q.Encode()

	return u
}

func hasDataUrl(id string, mt gog_integration.Media, pts ...vangogh_local_data.ProductType) *url.URL {
	u := &url.URL{
		Scheme: vangoghScheme,
		Host:   vangoghHost(),
		Path:   hasDataEndpoint,
	}

	productTypes := make([]string, 0, len(pts))
	for _, pt := range pts {
		productTypes = append(productTypes, pt.String())
	}

	q := u.Query()
	q.Set(vangogh_local_data.ProductTypeProperty, strings.Join(productTypes, ","))
	q.Set("media", mt.String())
	q.Set(vangogh_local_data.IdProperty, id)
	u.RawQuery = q.Encode()

	return u
}

func wishlistUrl(id string, mt gog_integration.Media) *url.URL {
	u := &url.URL{
		Scheme: vangoghScheme,
		Host:   vangoghHost(),
		Path:   wishlistEndpoint,
	}

	q := u.Query()
	q.Set("media", mt.String())
	q.Set(vangogh_local_data.IdProperty, id)
	u.RawQuery = q.Encode()

	return u
}

func tagUrl(id string, tags []string) *url.URL {
	u := &url.URL{
		Scheme: vangoghScheme,
		Host:   vangoghHost(),
		Path:   tagEndpoint,
	}

	q := u.Query()
	q.Set(vangogh_local_data.IdProperty, id)
	q.Set("tags", strings.Join(tags, ","))
	u.RawQuery = q.Encode()

	return u
}

func localTagUrl(id string, tags []string) *url.URL {
	u := &url.URL{
		Scheme: vangoghScheme,
		Host:   vangoghHost(),
		Path:   localTagEndpoint,
	}

	q := u.Query()
	q.Set(vangogh_local_data.IdProperty, id)
	q.Set("tags", strings.Join(tags, ","))
	u.RawQuery = q.Encode()

	return u
}

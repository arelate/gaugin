package api

import (
	"encoding/json"
	"fmt"
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_local_data"
	"html/template"
	"net/http"
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
)

type listProductViewModel struct {
	Id         string
	Title      string
	Developers []string
	Publisher  string
}

type listViewModel struct {
	Context  string
	Products []listProductViewModel
}

type productViewModel struct {
	Id string
	// text properties
	Title             string
	Image             string
	Tags              []string
	OperatingSystems  []string
	Rating            string
	Developers        []string
	Publisher         string
	Series            string
	Genres            []string
	Features          []string
	LanguageCodes     []string
	GlobalReleaseDate string
	GOGReleaseDate    string
	GOGOrderDate      string
	IncludesGames     []string
	IsIncludedByGames []string
	RequiresGames     []string
	IsRequiredByGames []string
	Types             []string
	// screenshots
	Screenshots []string
	// video-ids
	Videos []string
	// downloads
	Downloads vangogh_local_data.DownloadsList
}

func propertyFromRedux(redux map[string][]string, property string) string {
	properties := propertiesFromRedux(redux, property)
	if len(properties) > 0 {
		return properties[0]
	}
	return ""
}

func propertiesFromRedux(redux map[string][]string, property string) []string {
	if val, ok := redux[property]; ok {
		return val
	} else {
		return []string{}
	}
}

func listViewModelFromRedux(order []string, redux map[string]map[string][]string) *listViewModel {
	lvm := &listViewModel{
		Products: make([]listProductViewModel, 0, len(redux)),
	}
	for _, id := range order {
		properties, ok := redux[id]
		if !ok {
			continue
		}
		lvm.Products = append(
			lvm.Products,
			listProductViewModel{
				Id:         id,
				Title:      propertyFromRedux(properties, "title"),
				Developers: propertiesFromRedux(properties, "developers"),
				Publisher:  propertyFromRedux(properties, "publisher"),
			})
	}
	return lvm
}

func productViewModelFromRedux(redux map[string]map[string][]string) (*productViewModel, error) {
	switch len(redux) {
	case 0:
		return nil, fmt.Errorf("empty rdx")
	case 1:
		for id, rdx := range redux {

			return &productViewModel{
				Id:                id,
				Image:             propertyFromRedux(rdx, vangogh_local_data.ImageProperty),
				Title:             propertyFromRedux(rdx, vangogh_local_data.TitleProperty),
				Tags:              propertiesFromRedux(rdx, vangogh_local_data.TagIdProperty),
				OperatingSystems:  propertiesFromRedux(rdx, vangogh_local_data.OperatingSystemsProperty),
				Rating:            propertyFromRedux(rdx, vangogh_local_data.RatingProperty),
				Developers:        propertiesFromRedux(rdx, vangogh_local_data.DevelopersProperty),
				Publisher:         propertyFromRedux(rdx, vangogh_local_data.PublisherProperty),
				Series:            propertyFromRedux(rdx, vangogh_local_data.SeriesProperty),
				Genres:            propertiesFromRedux(rdx, vangogh_local_data.GenresProperty),
				Features:          propertiesFromRedux(rdx, vangogh_local_data.FeaturesProperty),
				LanguageCodes:     propertiesFromRedux(rdx, vangogh_local_data.LanguageCodeProperty),
				GlobalReleaseDate: propertyFromRedux(rdx, vangogh_local_data.GlobalReleaseDateProperty),
				GOGReleaseDate:    propertyFromRedux(rdx, vangogh_local_data.GOGReleaseDateProperty),
				GOGOrderDate:      propertyFromRedux(rdx, vangogh_local_data.GOGOrderDateProperty),
				IncludesGames:     propertiesFromRedux(rdx, vangogh_local_data.IncludesGamesProperty),
				IsIncludedByGames: propertiesFromRedux(rdx, vangogh_local_data.IsIncludedByGamesProperty),
				RequiresGames:     propertiesFromRedux(rdx, vangogh_local_data.RequiresGamesProperty),
				IsRequiredByGames: propertiesFromRedux(rdx, vangogh_local_data.IsRequiredByGamesProperty),
				Types:             propertiesFromRedux(rdx, vangogh_local_data.TypesProperty),
				Screenshots:       propertiesFromRedux(rdx, vangogh_local_data.ScreenshotsProperty),
				Videos:            propertiesFromRedux(rdx, vangogh_local_data.VideoIdProperty),
			}, nil
		}
	default:
		return nil, fmt.Errorf("too many ids, rdx")
	}
	return nil, nil
}

func funcMap() template.FuncMap {
	return template.FuncMap{
		"productTitle":  productTitle,
		"productId":     productId,
		"downloadTitle": downloadTitle,
		"formatBytes":   formatBytes,
	}
}

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

const titleIdSep = " ("

func productTitle(s string) string {
	if strings.Contains(s, titleIdSep) {
		return s[:strings.LastIndex(s, titleIdSep)]
	}
	return s
}

func productId(s string) string {
	if strings.Contains(s, titleIdSep) {
		return s[strings.LastIndex(s, titleIdSep)+len(titleIdSep) : len(s)-1]
	}
	return ""
}

func downloadTitle(d vangogh_local_data.Download) string {
	return d.String()
}

func formatBytes(b int) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}

func getKeys(client *http.Client, pt vangogh_local_data.ProductType, mt gog_integration.Media) ([]string, error) {
	ku := keysUrl(pt, mt)
	resp, err := client.Get(ku.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var keys []string
	err = json.NewDecoder(resp.Body).Decode(&keys)
	return keys, err
}

func getDownloads(
	client *http.Client,
	id string,
	operatingSystems []vangogh_local_data.OperatingSystem,
	languageCodes []string) (vangogh_local_data.DownloadsList, error) {
	du := downloadsUrl(id, operatingSystems, languageCodes)
	resp, err := client.Get(du.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var dl vangogh_local_data.DownloadsList
	err = json.NewDecoder(resp.Body).Decode(&dl)
	return dl, err
}

func getSearch(client *http.Client, q url.Values) ([]string, error) {
	su := searchUrl(q)
	resp, err := client.Get(su.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var keys []string
	err = json.NewDecoder(resp.Body).Decode(&keys)
	return keys, err
}

func getAllRedux(
	client *http.Client,
	pt vangogh_local_data.ProductType,
	mt gog_integration.Media, properties ...string) (map[string]map[string][]string, error) {
	ru := allReduxUrl(pt, mt, properties...)
	resp, err := client.Get(ru.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var redux map[string]map[string][]string
	err = json.NewDecoder(resp.Body).Decode(&redux)
	return redux, err
}

func getRedux(
	client *http.Client,
	id string,
	properties ...string) (map[string]map[string][]string, error) {
	ru := reduxUrl(id, properties...)
	resp, err := client.Get(ru.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var redux map[string]map[string][]string
	err = json.NewDecoder(resp.Body).Decode(&redux)
	return redux, err
}

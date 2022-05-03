package rest

import (
	"encoding/gob"
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_local_data"
	"net/http"
	"net/url"
	"time"
)

const (
	lastModified          = "Last-Modified"
	ifModifiedSinceHeader = "If-Modified-Since"
)

var (
	urlLastModified = make(map[string]int64)
	searchCache     = make(map[string][]string)
	digestsCache    = make(map[string]map[string][]string)
	reduxCache      = make(map[string]map[string]map[string][]string)
	downloadsCache  = make(map[string]vangogh_local_data.DownloadsList)
)

func getUpdates(
	client *http.Client,
	mt gog_integration.Media,
	since int) (map[string][]string, error) {
	uu := updatesUrl(mt, since)
	resp, err := client.Get(uu.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var updates map[string][]string
	err = gob.NewDecoder(resp.Body).Decode(&updates)
	return updates, err
}

func getThroughCache[T any](client *http.Client, u *url.URL, cache map[string]T) (T, error) {

	var data T

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return data, err
	}

	if lmt, ok := urlLastModified[u.String()]; ok {
		req.Header.Set(ifModifiedSinceHeader, time.Unix(lmt, 0).UTC().Format(time.RFC1123))
	}

	resp, err := client.Do(req)
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()

	if lm := resp.Header.Get(lastModified); lm != "" {
		if lmt, err := time.Parse(time.RFC1123, lm); err != nil {
			return data, err
		} else {
			urlLastModified[u.String()] = lmt.UTC().Unix()
		}
	}

	if resp.StatusCode == http.StatusNotModified {
		if data, ok := cache[u.String()]; ok {
			return data, nil
		}
	}

	err = gob.NewDecoder(resp.Body).Decode(&data)

	cache[u.String()] = data

	return data, err
}

func getDownloads(
	client *http.Client,
	id string,
	operatingSystems []vangogh_local_data.OperatingSystem,
	languageCodes []string) (vangogh_local_data.DownloadsList, error) {
	return getThroughCache(client, downloadsUrl(id, operatingSystems, languageCodes), downloadsCache)
}

func getRedux(
	client *http.Client,
	id string,
	properties ...string) (map[string]map[string][]string, error) {
	return getThroughCache(client, reduxUrl(id, properties...), reduxCache)
}

func getSearch(client *http.Client, q url.Values) ([]string, error) {
	return getThroughCache(client, searchUrl(q), searchCache)
}

func getDigests(client *http.Client,
	properties ...string) (map[string][]string, error) {
	return getThroughCache(client, digestUrl(properties...), digestsCache)
}

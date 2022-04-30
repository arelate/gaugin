package rest

import (
	"encoding/gob"
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_local_data"
	"net/http"
	"net/url"
	"time"
)

var urlLastModified = make(map[string]int64)
var searchKeysCache = make(map[string][]string)

func getKeys(client *http.Client, pt vangogh_local_data.ProductType, mt gog_integration.Media, count int) ([]string, error) {
	ku := keysUrl(pt, mt, count)
	resp, err := client.Get(ku.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var keys []string
	err = gob.NewDecoder(resp.Body).Decode(&keys)
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
	err = gob.NewDecoder(resp.Body).Decode(&dl)
	return dl, err
}

func getSearch(client *http.Client, q url.Values) ([]string, error) {
	su := searchUrl(q)
	req, err := http.NewRequest(http.MethodGet, su.String(), nil)
	if err != nil {
		return nil, err
	}
	if lmt, ok := urlLastModified[su.String()]; ok {
		req.Header.Set("If-Modified-Since", time.Unix(lmt, 0).UTC().Format(time.RFC1123))
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if lm := resp.Header.Get("Last-Modified"); lm != "" {
		if lmt, err := time.Parse(time.RFC1123, lm); err != nil {
			return nil, err
		} else {
			urlLastModified[su.String()] = lmt.UTC().Unix()
		}
	}

	if resp.StatusCode == http.StatusNotModified {
		return searchKeysCache[su.String()], nil
	}

	var keys []string
	err = gob.NewDecoder(resp.Body).Decode(&keys)

	searchKeysCache[su.String()] = keys

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
	err = gob.NewDecoder(resp.Body).Decode(&redux)
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
	err = gob.NewDecoder(resp.Body).Decode(&redux)
	return redux, err
}

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

func getDigests(client *http.Client,
	properties ...string) (map[string][]string, error) {
	du := digestUrl(properties...)
	resp, err := client.Get(du.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var digests map[string][]string
	err = gob.NewDecoder(resp.Body).Decode(&digests)
	return digests, err
}

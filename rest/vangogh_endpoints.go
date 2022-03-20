package rest

import (
	"encoding/gob"
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_local_data"
	"net/http"
	"net/url"
)

func getKeys(client *http.Client, pt vangogh_local_data.ProductType, mt gog_integration.Media) ([]string, error) {
	ku := keysUrl(pt, mt)
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
	resp, err := client.Get(su.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var keys []string
	err = gob.NewDecoder(resp.Body).Decode(&keys)
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

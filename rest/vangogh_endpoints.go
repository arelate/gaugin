package rest

import (
	"encoding/gob"
	"encoding/json"
	"github.com/arelate/gog_integration"
	"github.com/arelate/steam_integration"
	"github.com/arelate/vangogh_local_data"
	"net/http"
	"net/url"
	"sync"
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
	mtx             = sync.Mutex{}
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
			mtx.Lock()
			urlLastModified[u.String()] = lmt.UTC().Unix()
			mtx.Unlock()
		}
	}

	if resp.StatusCode == http.StatusNotModified {
		if data, ok := cache[u.String()]; ok {
			return data, nil
		}
	}

	err = gob.NewDecoder(resp.Body).Decode(&data)

	mtx.Lock()
	cache[u.String()] = data
	mtx.Unlock()

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

func getHasRedux(
	client *http.Client,
	id string,
	properties ...string) (map[string]map[string][]string, error) {
	return getThroughCache(client, hasReduxUrl(id, properties...), reduxCache)
}

func getSearch(client *http.Client, q url.Values) ([]string, error) {
	return getThroughCache(client, searchUrl(q), searchCache)
}

func getDigests(client *http.Client, properties ...string) (map[string][]string, error) {
	return getThroughCache(client, digestUrl(properties...), digestsCache)
}

func getSteamAppNews(client *http.Client, id string) (*steam_integration.AppNews, error) {
	sanu := steamAppNewsUrl(id)
	resp, err := client.Get(sanu.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data map[string]steam_integration.GetNewsForAppResponse
	err = json.NewDecoder(resp.Body).Decode(&data)

	if steamAppNews, ok := data[id]; ok {
		return &steamAppNews.AppNews, nil
	}

	return nil, err
}

func getHasData(
	client *http.Client,
	id string,
	mt gog_integration.Media,
	pts ...vangogh_local_data.ProductType) (map[vangogh_local_data.ProductType]bool, error) {

	hasData := make(map[vangogh_local_data.ProductType]bool)

	hdu := hasDataUrl(id, mt, pts...)
	resp, err := client.Get(hdu.String())
	if err != nil {
		return hasData, err
	}
	defer resp.Body.Close()

	var data map[string]map[string]string
	err = gob.NewDecoder(resp.Body).Decode(&data)

	for _, pt := range pts {
		if hasProductType, ok := data[pt.String()]; ok {
			if hd, sure := hasProductType[id]; sure {
				hasData[pt] = hd == "true"
			}
		}
	}

	return hasData, err
}

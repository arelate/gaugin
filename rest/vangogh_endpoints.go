package rest

import (
	"encoding/gob"
	"fmt"
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
	hasDataCache    = make(map[string]map[string]map[string]string)
	dataCache       = make(map[string]map[string]interface{})
	searchCache     = make(map[string][]string)
	digestsCache    = make(map[string]map[string][]string)
	hasReduxCache   = make(map[string]map[string]map[string][]string)
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
		if cdata, ok := cache[u.String()]; ok {
			return cdata, nil
		}
	}

	if err = gob.NewDecoder(resp.Body).Decode(&data); err == nil {
		mtx.Lock()
		cache[u.String()] = data
		mtx.Unlock()
	}

	return data, err
}

func getDownloads(
	client *http.Client,
	id string,
	operatingSystems []vangogh_local_data.OperatingSystem,
	languageCodes []string) (vangogh_local_data.DownloadsList, error) {
	return getThroughCache(client, downloadsUrl(id, operatingSystems, languageCodes), downloadsCache)
}

func getHasRedux(
	client *http.Client,
	id string,
	properties ...string) (map[string]map[string][]string, error) {
	return getThroughCache(client, hasReduxUrl(id, properties...), hasReduxCache)
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

func getDigests(client *http.Client, properties ...string) (map[string][]string, error) {
	return getThroughCache(client, digestUrl(properties...), digestsCache)
}

func getHasData(
	client *http.Client,
	id string,
	mt gog_integration.Media,
	pts ...vangogh_local_data.ProductType) (map[string]map[string]string, error) {
	return getThroughCache(client, hasDataUrl(id, mt, pts...), hasDataCache)
}

func getData(
	client *http.Client,
	id string,
	pt vangogh_local_data.ProductType,
	mt gog_integration.Media) (map[string]interface{}, error) {
	return getThroughCache(client, dataUrl(id, pt, mt), dataCache)
}

func getSteamAppNews(client *http.Client, id string) (*steam_integration.AppNews, error) {

	data, err := getData(client, id, vangogh_local_data.SteamAppNews, gog_integration.Game)
	if err != nil {
		return nil, err
	}

	if getNewsForAppResponseData, ok := data[id]; ok {
		if getNewsForAppResponse, sure := getNewsForAppResponseData.(steam_integration.GetNewsForAppResponse); sure {
			appNews := getNewsForAppResponse.AppNews
			return &appNews, nil
		}
	}

	return nil, err
}

func wishlistMethod(client *http.Client, method string, id string, mt gog_integration.Media) error {
	wu := wishlistUrl(id, mt)

	req, err := http.NewRequest(method, wu.String(), nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	return nil
}

func putWishlist(client *http.Client, id string, mt gog_integration.Media) error {
	return wishlistMethod(client, http.MethodPut, id, mt)
}

func deleteWishlist(client *http.Client, id string, mt gog_integration.Media) error {
	return wishlistMethod(client, http.MethodDelete, id, mt)
}

func patchTag(client *http.Client, id string, tags []string) error {
	tu := tagUrl(id, tags)

	req, err := http.NewRequest(http.MethodPatch, tu.String(), nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	return nil
}

func patchLocalTag(client *http.Client, id string, tags []string) error {
	ltu := localTagUrl(id, tags)

	req, err := http.NewRequest(http.MethodPatch, ltu.String(), nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	return nil
}

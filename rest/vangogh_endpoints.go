package rest

import (
	"encoding/gob"
	"fmt"
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
	digestsCache    = make(vangogh_local_data.IdReduxAssets)
	hasReduxCache   = make(map[string]vangogh_local_data.IdReduxAssets)
	reduxCache      = make(map[string]vangogh_local_data.IdReduxAssets)
	downloadsCache  = make(map[string]vangogh_local_data.DownloadsList)
	mtx             = sync.Mutex{}
)

func getThroughCache[T any](client *http.Client, u *url.URL, cache map[string]T) (T, bool, error) {

	var data T

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return data, false, err
	}

	if lmt, ok := urlLastModified[u.String()]; ok {
		req.Header.Set(ifModifiedSinceHeader, time.Unix(lmt, 0).UTC().Format(time.RFC1123))
	}

	resp, err := client.Do(req)
	if err != nil {
		return data, false, err
	}
	defer resp.Body.Close()

	if lm := resp.Header.Get(lastModified); lm != "" {
		if lmt, err := time.Parse(time.RFC1123, lm); err != nil {
			return data, false, err
		} else {
			mtx.Lock()
			urlLastModified[u.String()] = lmt.UTC().Unix()
			mtx.Unlock()
		}
	}

	if resp.StatusCode == http.StatusNotModified {
		if cdata, ok := cache[u.String()]; ok {
			return cdata, true, nil
		}
	}

	if err = gob.NewDecoder(resp.Body).Decode(&data); err == nil {
		mtx.Lock()
		cache[u.String()] = data
		mtx.Unlock()
	}

	return data, false, err
}

func getDownloads(
	client *http.Client,
	id string,
	operatingSystems []vangogh_local_data.OperatingSystem,
	languageCodes []string) (vangogh_local_data.DownloadsList, bool, error) {
	return getThroughCache(client, downloadsUrl(id, operatingSystems, languageCodes), downloadsCache)
}

func getHasRedux(
	client *http.Client,
	id string,
	properties ...string) (vangogh_local_data.IdReduxAssets, bool, error) {
	return getThroughCache(client, hasReduxUrl(id, properties...), hasReduxCache)
}

func getRedux(
	client *http.Client,
	id string,
	all bool,
	properties ...string) (vangogh_local_data.IdReduxAssets, bool, error) {
	if all && len(properties) > 1 {
		return nil, false, fmt.Errorf("cannot use all with more than 1 property")
	}
	return getThroughCache(client, reduxUrl(id, all, properties...), reduxCache)
}

func getSearch(client *http.Client, q url.Values) ([]string, bool, error) {
	return getThroughCache(client, searchUrl(q), searchCache)
}

func getDigests(client *http.Client, properties ...string) (map[string][]string, bool, error) {
	return getThroughCache(client, digestUrl(properties...), digestsCache)
}

func getHasData(
	client *http.Client,
	id string,
	pts ...vangogh_local_data.ProductType) (map[string]map[string]string, bool, error) {
	return getThroughCache(client, hasDataUrl(id, pts...), hasDataCache)
}

func getData(
	client *http.Client,
	id string,
	pt vangogh_local_data.ProductType) (map[string]interface{}, bool, error) {
	return getThroughCache(client, dataUrl(id, pt), dataCache)
}

func getSteamNews(client *http.Client, id string) (*steam_integration.AppNews, bool, error) {

	data, cached, err := getData(client, id, vangogh_local_data.SteamAppNews)
	if err != nil {
		return nil, cached, err
	}

	if getNewsForAppResponseData, ok := data[id]; ok {
		if getNewsForAppResponse, sure := getNewsForAppResponseData.(steam_integration.GetNewsForAppResponse); sure {
			appNews := getNewsForAppResponse.AppNews
			return &appNews, cached, nil
		}
	}

	return nil, cached, err
}

func getSteamReviews(client *http.Client, id string) (*steam_integration.AppReviews, bool, error) {
	data, cached, err := getData(client, id, vangogh_local_data.SteamReviews)
	if err != nil {
		return nil, cached, err
	}

	if appReviewsData, ok := data[id]; ok {
		if appReviews, sure := appReviewsData.(steam_integration.AppReviews); sure {
			return &appReviews, cached, nil
		}
	}

	return nil, cached, err

}

func wishlistMethod(client *http.Client, method string, id string) error {
	wu := wishlistUrl(id)

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

func putWishlist(client *http.Client, id string) error {
	return wishlistMethod(client, http.MethodPut, id)
}

func deleteWishlist(client *http.Client, id string) error {
	return wishlistMethod(client, http.MethodDelete, id)
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

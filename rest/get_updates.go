package rest

import (
	"errors"
	"github.com/arelate/gaugin/stencil_app"
	"github.com/boggydigital/kvas"
	"golang.org/x/exp/slices"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/arelate/gaugin/gaugin_middleware"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"golang.org/x/exp/maps"
)

func GetUpdates(w http.ResponseWriter, r *http.Request) {

	// GET /updates

	st := gaugin_middleware.NewServerTimings()

	start := time.Now()
	updRdx, cached, err := getRedux(
		http.DefaultClient,
		"",
		true,
		vangogh_local_data.LastSyncUpdatesProperty)

	if cached {
		st.SetFlag("updRdx-cached")
	}
	st.Set("updRdx", time.Since(start).Milliseconds())

	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	updates := make(map[string][]string)
	for section, rdx := range updRdx {
		ids := rdx[vangogh_local_data.LastSyncUpdatesProperty]
		if len(ids) > 0 {
			updates[section] = ids
		}
	}

	keys := make(map[string]bool)
	for _, ids := range updates {
		for _, id := range ids {
			keys[id] = true
		}
	}

	ids := maps.Keys(keys)
	sort.Strings(ids)

	start = time.Now()
	dataRdx, cached, err := getRedux(
		http.DefaultClient,
		strings.Join(ids, ","),
		false,
		stencil_app.ProductsProperties...)

	if cached {
		st.SetFlag("dataRdx-cached")
	}
	st.Set("dataRdx", time.Since(start).Milliseconds())

	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	start = time.Now()
	syncRdx, cached, err := getRedux(
		http.DefaultClient,
		vangogh_local_data.SyncCompleteKey,
		false,
		vangogh_local_data.SyncEventsProperty)

	if cached {
		st.SetFlag("syncRdx-cached")
	}
	st.Set("syncRdx", time.Since(start).Milliseconds())

	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	updated := "recently"
	syncDra := NewDataRdx(syncRdx)
	if scs, ok := syncDra.GetFirstVal(vangogh_local_data.SyncEventsProperty, vangogh_local_data.SyncCompleteKey); ok {
		if sci, err := strconv.ParseInt(scs, 10, 64); err == nil {
			updated = time.Unix(sci, 0).Format(time.RFC1123)
		}
	}

	dra := NewDataRdx(dataRdx)

	//uvm := view_models.NewUpdates(updates, dataRdx, syncRdx[vangogh_local_data.SyncCompleteKey])

	gaugin_middleware.DefaultHeaders(st, w)

	sections := maps.Keys(updates)
	sort.Strings(sections)

	var caser = cases.Title(language.Russian)

	sectionTitles := make(map[string]string)
	for t, _ := range updates {
		sectionTitles[t] = caser.String(t)
	}

	if err := app.RenderGroup(stencil_app.NavUpdates, sections, updates, sectionTitles, updated, dra, w); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	//if err := tmpl.ExecuteTemplate(w, "updates-page", uvm); err != nil {
	//	http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
	//	return
	//}
}

var DataReduxAssetsReadOnlyError = errors.New("data redux assets is read-only")

type DataReduxAssets struct {
	rdx     map[string]map[string][]string
	modTime int64
}

func NewDataRdx(data map[string]map[string][]string) *DataReduxAssets {
	return &DataReduxAssets{
		rdx:     data,
		modTime: time.Now().UTC().Unix(),
	}
}

func NewEmpty(properties []string) *DataReduxAssets {
	dra := NewDataRdx(make(map[string]map[string][]string))
	dra.rdx[""] = make(map[string][]string)
	for _, p := range properties {
		dra.rdx[""][p] = nil
	}
	return dra
}

func (dra *DataReduxAssets) Keys(asset string) []string {
	keys := make([]string, 0)
	for id, pvs := range dra.rdx {
		if _, ok := pvs[asset]; ok {
			keys = append(keys, id)
		}
	}
	return keys
}

func (dra *DataReduxAssets) Has(asset string) bool {
	for _, pvs := range dra.rdx {
		if _, ok := pvs[asset]; ok {
			return true
		}
	}
	return false
}

func (dra *DataReduxAssets) HasKey(asset, key string) bool {
	if pvs, ok := dra.rdx[key]; ok {
		if _, ok := pvs[asset]; ok {
			return true
		}
	}
	return false
}

func (dra *DataReduxAssets) HasVal(asset, key, val string) bool {
	if pvs, ok := dra.rdx[key]; ok {
		if vals, ok := pvs[asset]; ok {
			return slices.Contains(vals, val)
		}
	}
	return false
}

func (dra *DataReduxAssets) AddVal(asset, key, val string) error {
	return DataReduxAssetsReadOnlyError
}

func (dra *DataReduxAssets) ReplaceValues(asset, key string, values ...string) error {
	return DataReduxAssetsReadOnlyError
}

func (dra *DataReduxAssets) BatchReplaceValues(asset string, keyValues map[string][]string) error {
	return DataReduxAssetsReadOnlyError
}

func (dra *DataReduxAssets) CutVal(asset, key, val string) error {
	return DataReduxAssetsReadOnlyError
}

func (dra *DataReduxAssets) GetAllValues(asset, key string) ([]string, bool) {
	if pvs, ok := dra.rdx[key]; ok {
		if vals, ok := pvs[asset]; ok {
			return vals, true
		}
	}
	return nil, false
}

func (dra *DataReduxAssets) GetAllUnchangedValues(asset, key string) ([]string, bool) {
	return dra.GetAllValues(asset, key)
}

func (dra *DataReduxAssets) GetFirstVal(asset, key string) (string, bool) {
	if vals, ok := dra.GetAllValues(asset, key); ok {
		if len(vals) > 0 {
			return vals[0], true
		}
	}
	return "", false
}

func (dra *DataReduxAssets) IsSupported(assets ...string) error {
	for _, a := range assets {
		if !dra.Has(a) {
			return errors.New("unsupported asset " + a)
		}
	}
	return nil
}

func (dra *DataReduxAssets) Match(query map[string][]string, anyCase bool) map[string]bool {
	//FIXME
	return nil
}

func (dra *DataReduxAssets) RefreshReduxAssets() (kvas.ReduxAssets, error) {
	return dra, nil
}

func (dra *DataReduxAssets) ReduxAssetsModTime() (int64, error) {
	return dra.modTime, nil
}

func (dra *DataReduxAssets) Sort(ids []string, desc bool, sortBy ...string) ([]string, error) {
	//FIXME
	return ids, nil
}

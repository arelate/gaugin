package rest

import (
	"errors"
	"github.com/boggydigital/kvas"
	"golang.org/x/exp/slices"
	"io"
	"time"
)

type IdReduxAssets = map[string]map[string][]string

type IRAProxy struct {
	rdx     IdReduxAssets
	modTime int64
}

func NewIRAProxy(data IdReduxAssets) *IRAProxy {
	return &IRAProxy{
		rdx:     data,
		modTime: time.Now().UTC().Unix(),
	}
}

func NewEmptyIRAProxy(properties []string) *IRAProxy {
	dra := NewIRAProxy(make(IdReduxAssets))
	dra.rdx[""] = make(map[string][]string)
	for _, p := range properties {
		dra.rdx[""][p] = nil
	}
	return dra
}

func (irap *IRAProxy) Keys(asset string) []string {
	keys := make([]string, 0)
	for id, pvs := range irap.rdx {
		if _, ok := pvs[asset]; ok {
			keys = append(keys, id)
		}
	}
	return keys
}

func (irap *IRAProxy) HasAsset(asset string) bool {
	for _, pvs := range irap.rdx {
		if _, ok := pvs[asset]; ok {
			return true
		}
	}
	return false
}

func (irap *IRAProxy) HasKey(asset, key string) bool {
	if pvs, ok := irap.rdx[key]; ok {
		if _, ok := pvs[asset]; ok {
			return true
		}
	}
	return false
}

func (irap *IRAProxy) HasValue(asset, key, val string) bool {
	if pvs, ok := irap.rdx[key]; ok {
		if vals, ok := pvs[asset]; ok {
			return slices.Contains(vals, val)
		}
	}
	return false
}

func (irap *IRAProxy) MustHave(assets ...string) error {
	for _, a := range assets {
		if len(irap.rdx) > 0 && !irap.HasAsset(a) {
			return errors.New("unsupported asset " + a)
		}
	}
	return nil
}

func (irap *IRAProxy) GetAllValues(asset, key string) ([]string, bool) {
	if pvs, ok := irap.rdx[key]; ok {
		if vals, ok := pvs[asset]; ok {
			return vals, true
		}
	}
	return nil, false
}

func (irap *IRAProxy) GetFirstVal(asset, key string) (string, bool) {
	if vals, ok := irap.GetAllValues(asset, key); ok {
		if len(vals) > 0 {
			return vals[0], true
		}
	}
	return "", false
}

func (irap *IRAProxy) MatchAsset(asset string, terms []string, scope []string, options ...kvas.MatchOption) []string {
	return nil
}

func (irap *IRAProxy) Match(query map[string][]string, options ...kvas.MatchOption) []string {
	//FIXME
	return nil
}

func (irap *IRAProxy) Refresh() error {
	return nil
}

func (irap *IRAProxy) ModTime() (int64, error) {
	return irap.modTime, nil
}

func (irap *IRAProxy) Sort(ids []string, desc bool, sortBy ...string) ([]string, error) {
	return ids, errors.New("IRAProxy doesn't support sorting")
}

func (irap *IRAProxy) RefreshReader() (kvas.ReadableRedux, error) {
	return irap, nil
}

func (irap *IRAProxy) Export(w io.Writer, ids ...string) error {
	//FIXME
	return nil
}

func (irap *IRAProxy) Merge(idPropertyValues map[string]map[string][]string) {
	for id, pv := range idPropertyValues {
		if irap.rdx[id] == nil {
			irap.rdx[id] = make(map[string][]string)
		}
		for p, v := range pv {
			irap.rdx[id][p] = append(irap.rdx[id][p], v...)
		}
	}
}

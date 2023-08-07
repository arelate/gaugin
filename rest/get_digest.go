package rest

import (
	"encoding/json"
	"github.com/arelate/gaugin/gaugin_middleware"
	"github.com/arelate/gaugin/stencil_app"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/http"
	"time"
)

func GetDigest(w http.ResponseWriter, r *http.Request) {

	// GET /digest?property

	property := r.URL.Query().Get("property")
	if property == "" {
		http.Error(w, nod.ErrorStr("missing digest property"), http.StatusBadRequest)
		return
	}

	start := time.Now()
	st := gaugin_middleware.NewServerTimings()

	var values []string
	valueTitles := make(map[string]string)
	var digests map[string][]string
	var cached bool
	var err error

	switch property {
	case vangogh_local_data.SortProperty:
		values = []string{
			vangogh_local_data.GlobalReleaseDateProperty,
			vangogh_local_data.GOGReleaseDateProperty,
			vangogh_local_data.GOGOrderDateProperty,
			vangogh_local_data.TitleProperty,
			vangogh_local_data.RatingProperty,
			vangogh_local_data.DiscountPercentageProperty,
			vangogh_local_data.HLTBHoursToCompleteMainProperty,
			vangogh_local_data.HLTBHoursToCompletePlusProperty,
			vangogh_local_data.HLTBHoursToComplete100Property}
	case vangogh_local_data.DescendingProperty:
		values = []string{
			vangogh_local_data.TrueValue,
			vangogh_local_data.FalseValue}
	case vangogh_local_data.TagIdProperty:
		tagNamesRedux, _, err := getRedux(http.DefaultClient, "", true, vangogh_local_data.TagNameProperty)
		if err != nil {
			http.Error(w, nod.ErrorStr("missing digest property"), http.StatusBadRequest)
			return
		}
		for tagId, tagNames := range tagNamesRedux {
			if tns, ok := tagNames[vangogh_local_data.TagNameProperty]; ok && len(tns) > 0 {
				valueTitles[tagId] = tns[0]
			}
		}
	default:
		digests, cached, err = getDigests(http.DefaultClient, property)
	}

	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	if len(digests) > 0 {
		values = digests[property]
	}

	if cached {
		st.SetFlag("getDigests-cached")
	}
	st.Set("getDigests", time.Since(start).Milliseconds())

	gaugin_middleware.DefaultHeaders(st, w)

	addedValueTitles, err := addTitles(property, values)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	for v, t := range addedValueTitles {
		valueTitles[v] = t
	}

	if err := json.NewEncoder(w).Encode(valueTitles); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}
}

func addTitles(property string, values []string) (map[string]string, error) {
	valueTitles := make(map[string]string)

	switch property {
	case vangogh_local_data.TagIdProperty:
		// do nothing already filled earlier
	default:
		for _, v := range values {
			if title, ok := stencil_app.PropertyTitles[v]; ok {
				valueTitles[v] = title
			} else {
				valueTitles[v] = v
			}
		}
	}
	return valueTitles, nil
}

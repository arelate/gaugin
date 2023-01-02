package rest

import (
	"encoding/json"
	"github.com/arelate/gaugin/gaugin_middleware"
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

	var digests map[string][]string
	var cached bool
	var err error

	switch property {
	case vangogh_local_data.SortProperty:
		digests, cached = map[string][]string{vangogh_local_data.SortProperty: {
			vangogh_local_data.GlobalReleaseDateProperty,
			vangogh_local_data.GOGReleaseDateProperty,
			vangogh_local_data.GOGOrderDateProperty,
			vangogh_local_data.TitleProperty,
			vangogh_local_data.RatingProperty,
			vangogh_local_data.DiscountPercentageProperty,
			vangogh_local_data.HLTBHoursToCompleteMainProperty,
			vangogh_local_data.HLTBHoursToCompletePlusProperty,
			vangogh_local_data.HLTBHoursToComplete100Property}}, true
	case vangogh_local_data.DescendingProperty:
		digests, cached = map[string][]string{vangogh_local_data.DescendingProperty: {
			vangogh_local_data.TrueValue,
			vangogh_local_data.FalseValue}}, true
	default:
		digests, cached, err = getDigests(http.DefaultClient, property)
	}

	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	if cached {
		st.SetFlag("getDigests-cached")
	}
	st.Set("getDigests", time.Since(start).Milliseconds())

	gaugin_middleware.DefaultHeaders(st, w)

	if err := json.NewEncoder(w).Encode(digests[property]); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}
}

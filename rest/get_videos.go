package rest

import (
	"github.com/arelate/gaugin/rest/compton_pages"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/kevlar"
	"github.com/boggydigital/nod"
	"net/http"
	"slices"
)

func GetVideos(w http.ResponseWriter, r *http.Request) {

	// GET /videos?id

	id := r.URL.Query().Get("id")

	idRedux, err := getRedux(
		http.DefaultClient,
		id,
		false,
		vangogh_local_data.VideoIdProperty)

	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	rdx := kevlar.ReduxProxy(idRedux)

	var videoIds []string
	if vids, ok := rdx.GetAllValues(vangogh_local_data.VideoIdProperty, id); ok {
		videoIds = vids
	}
	slices.Sort(videoIds)

	videoTitles := make(map[string]string)
	videoDurations := make(map[string]string)

	for _, vid := range videoIds {
		vpRedux, err := getRedux(
			http.DefaultClient,
			vid,
			false,
			vangogh_local_data.VideoTitleProperty,
			vangogh_local_data.VideoDurationProperty)

		if err != nil {
			http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
			return
		}

		vpRdx := kevlar.ReduxProxy(vpRedux)

		if vtp, ok := vpRdx.GetLastVal(vangogh_local_data.VideoTitleProperty, vid); ok {
			videoTitles[vid] = vtp
		}
		if vdp, ok := vpRdx.GetLastVal(vangogh_local_data.VideoDurationProperty, vid); ok {
			videoDurations[vid] = vdp
		}
	}

	p := compton_pages.Videos(videoIds, videoTitles, videoDurations)

	if err := p.WriteContent(w); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
	}
}

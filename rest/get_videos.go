package rest

import (
	"github.com/arelate/gaugin/rest/compton_data"
	"github.com/arelate/gaugin/rest/gaugin_styles"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/compton/consts/color"
	"github.com/boggydigital/compton/consts/direction"
	"github.com/boggydigital/compton/elements/details_summary"
	"github.com/boggydigital/compton/elements/els"
	"github.com/boggydigital/compton/elements/flex_items"
	"github.com/boggydigital/compton/elements/fspan"
	"github.com/boggydigital/compton/elements/iframe_expand"
	"github.com/boggydigital/compton/elements/recipes"
	"github.com/boggydigital/kevlar"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet_urls/youtube_urls"
	"net/http"
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

	mvRedux, err := getRedux(
		http.DefaultClient,
		"",
		true,
		vangogh_local_data.MissingVideoUrlProperty)

	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	rdx := kevlar.ReduxProxy(MergeIdPropertyValues(idRedux, mvRedux))

	videoIds := make([]string, 0)
	if vp, ok := rdx.GetAllValues(vangogh_local_data.VideoIdProperty, id); ok && len(vp) > 0 {
		videoIds = vp
	}

	localVideos := make([]string, 0, len(videoIds))

	for _, videoId := range videoIds {
		if !rdx.HasKey(vangogh_local_data.MissingVideoUrlProperty, videoId) {
			localVideos = append(localVideos, videoId)
		}
	}

	section := compton_data.VideosSection
	ifc := iframe_expand.IframeExpandContent(section, compton_data.SectionTitles[section]).
		AppendStyle(gaugin_styles.VideosStyle)

	pageStack := flex_items.FlexItems(ifc, direction.Column)
	ifc.Append(pageStack)

	if len(videoIds) == 0 {
		fs := fspan.Text(ifc, "Videos are not available for this product").
			ForegroundColor(color.Subtle)
		pageStack.Append(recipes.Center(ifc, fs))
	}

	for _, videoId := range localVideos {
		src := "/video?id=" + videoId
		posterSrc := "/thumbnails?id=" + videoId
		video := els.Video(src)
		video.SetAttribute("controls", "")
		video.SetAttribute("poster", posterSrc)
		pageStack.Append(video)
	}

	if len(videoIds) > 0 {

		detailsSummary := details_summary.
			Closed(ifc, "Watch on YouTube").
			ForegroundColor(color.Cyan)

		dsStack := flex_items.FlexItems(ifc, direction.Column)
		detailsSummary.Append(dsStack)

		for _, videoId := range videoIds {
			posterSrc := "/thumbnails?id=" + videoId
			originSrc := youtube_urls.VideoUrl(videoId).String()
			link := els.A(originSrc)
			link.SetAttribute("target", "_top")
			img := els.ImageLazy(posterSrc)
			img.SetAttribute("alt", "VideoId: "+videoId)
			link.Append(img)
			dsStack.Append(link)
		}

		pageStack.Append(detailsSummary)
	}

	if err := ifc.WriteContent(w); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
	}
}

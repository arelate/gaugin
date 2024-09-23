package rest

import (
	"github.com/arelate/gaugin/rest/compton_data"
	"github.com/arelate/gaugin/rest/gaugin_styles"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/color"
	"github.com/boggydigital/compton/consts/direction"
	"github.com/boggydigital/compton/elements/els"
	"github.com/boggydigital/compton/elements/flex_items"
	"github.com/boggydigital/compton/elements/fspan"
	"github.com/boggydigital/compton/elements/iframe_expand"
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

	rdx := kevlar.ReduxProxy(idRedux)

	videoIds := make([]string, 0)
	if vp, ok := rdx.GetAllValues(vangogh_local_data.VideoIdProperty, id); ok && len(vp) > 0 {
		videoIds = vp
	}

	section := compton_data.VideosSection
	ifc := iframe_expand.IframeExpandContent(section, compton_data.SectionTitles[section]).
		AppendStyle(gaugin_styles.VideosStyle)

	pageStack := flex_items.FlexItems(ifc, direction.Column)
	ifc.Append(pageStack)

	if len(videoIds) == 0 {
		fs := fspan.Text(ifc, "Videos are not available for this product").
			ForegroundColor(color.Subtle)
		pageStack.Append(flex_items.Center(ifc, fs))
	}

	for _, videoId := range videoIds {
		pageStack.Append(createVideo(ifc, videoId))
	}

	if err := ifc.WriteContent(w); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
	}
}

func createVideo(r compton.Registrar, videoId string) compton.Element {
	src := "/video?id=" + videoId
	posterSrc := "/thumbnails?id=" + videoId

	stack := flex_items.FlexItems(r, direction.Column)

	video := els.Video(src)
	video.SetAttribute("controls", "")
	video.SetAttribute("poster", posterSrc)
	stack.Append(video)

	originLink := els.AText("Watch at origin", youtube_urls.VideoUrl(videoId).String())
	originLink.AddClass("external")
	stack.Append(originLink)

	return stack
}

package rest

import (
	"github.com/arelate/gaugin/rest/compton_data"
	"github.com/arelate/gaugin/rest/gaugin_styles"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/compton/consts/direction"
	"github.com/boggydigital/compton/elements/els"
	"github.com/boggydigital/compton/elements/flex_items"
	"github.com/boggydigital/compton/elements/iframe_expand"
	"github.com/boggydigital/nod"
	"net/http"
)

func GetScreenshots(w http.ResponseWriter, r *http.Request) {

	// GET /screenshots?id

	id := r.URL.Query().Get("id")

	idRedux, err := getRedux(
		http.DefaultClient,
		id,
		false,
		vangogh_local_data.ScreenshotsProperty)

	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	var screenshots []string
	if rdx := idRedux[id]; rdx != nil {
		screenshots = rdx[vangogh_local_data.ScreenshotsProperty]
	}

	section := compton_data.ScreenshotsSection
	ifc := iframe_expand.IframeExpandContent(section, compton_data.SectionTitles[section]).
		AppendStyle(gaugin_styles.ScreenshotsStyle)

	pageStack := flex_items.FlexItems(ifc, direction.Column)
	ifc.Append(pageStack)

	for _, src := range screenshots {
		imageSrc := "/image?id=" + src
		link := els.A(imageSrc)
		link.SetAttribute("target", "_top")
		link.Append(els.Img(imageSrc))
		pageStack.Append(link)
	}

	if err := ifc.WriteContent(w); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
	}
}

package compton_pages

import (
	"github.com/arelate/gaugin/rest/compton_data"
	"github.com/arelate/gaugin/rest/compton_fragments"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/color"
	"github.com/boggydigital/compton/consts/direction"
	"github.com/boggydigital/compton/elements/flex_items"
	"github.com/boggydigital/compton/elements/fspan"
	"github.com/boggydigital/kevlar"
)

func Videos(id string, rdx kevlar.ReadableRedux) compton.Element {
	videoIds, _ := rdx.GetAllValues(vangogh_local_data.VideoIdProperty, id)

	s := compton_fragments.ProductSection(compton_data.VideosSection)

	pageStack := flex_items.FlexItems(s, direction.Column)
	s.Append(pageStack)

	if len(videoIds) == 0 {
		fs := fspan.Text(s, "Videos are not available for this product").
			ForegroundColor(color.Gray)
		pageStack.Append(flex_items.Center(s, fs))
	}

	for _, videoId := range videoIds {
		pageStack.Append(compton_fragments.Video(s, videoId))
	}

	return s
}

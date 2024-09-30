package compton_pages

import (
	"github.com/arelate/gaugin/rest/compton_data"
	"github.com/arelate/gaugin/rest/compton_fragments"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/direction"
	"github.com/boggydigital/compton/elements/flex_items"
	"github.com/boggydigital/kevlar"
)

func Downloads(id string, os vangogh_local_data.OperatingSystem, dls vangogh_local_data.DownloadsList, rdx kevlar.ReadableRedux) compton.Element {
	s := compton_fragments.ProductSection(compton_data.DownloadsSection)

	pageStack := flex_items.FlexItems(s, direction.Column)
	s.Append(pageStack)

	return s
}

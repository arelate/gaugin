package compton_pages

import (
	"github.com/arelate/gaugin/rest/compton_data"
	"github.com/arelate/gaugin/rest/compton_fragments"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/direction"
	"github.com/boggydigital/compton/elements/els"
	"github.com/boggydigital/compton/elements/flex_items"
	"github.com/boggydigital/kevlar"
)

var osOrder = []vangogh_local_data.OperatingSystem{
	vangogh_local_data.Windows,
	vangogh_local_data.MacOS,
	vangogh_local_data.Linux,
}

// Downloads will present available installers, DLCs in the following hierarchy:
// - Operating system heading - Installers and DLCs (separately)
// - title_values list of downloads by version
func Downloads(id string, clientOs vangogh_local_data.OperatingSystem, dls vangogh_local_data.DownloadsList, rdx kevlar.ReadableRedux) compton.Element {
	s := compton_fragments.ProductSection(compton_data.DownloadsSection)

	pageStack := flex_items.FlexItems(s, direction.Column)
	s.Append(pageStack)

	for ii, os := range osOrder {
		osHeading := els.H3Text(os.String())
		pageStack.Append(osHeading)
		if ii < len(osOrder)-1 {
			pageStack.Append(els.Hr())
		}

	}

	return s
}

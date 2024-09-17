package compton_fragments

import (
	"github.com/arelate/gaugin/rest/compton_data"
	"github.com/arelate/gaugin/rest/gaugin_styles"
	"github.com/boggydigital/compton/page"
)

func GauginPage(current string) *page.PageElement {
	return page.Page(
		PageTitle(current)).
		SetFavIconEmoji(compton_data.AppFavIconEmoji).
		SetCustomStyles(gaugin_styles.GauginStyle)
}

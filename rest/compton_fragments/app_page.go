package compton_fragments

import (
	"github.com/arelate/gaugin/rest/compton_data"
	"github.com/arelate/gaugin/rest/gaugin_styles"
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/direction"
	"github.com/boggydigital/compton/elements/flex_items"
	"github.com/boggydigital/compton/page"
)

func AppPage(current string) (p *page.PageElement, stack compton.Element) {
	p = page.Page(
		PageTitle(current)).
		SetFavIconEmoji(compton_data.AppFavIconEmoji).
		AppendStyle(gaugin_styles.AppStyle)

	stack = flex_items.FlexItems(p, direction.Column)
	p.Append(stack)

	return p, stack
}

package compton_pages

import (
	"github.com/arelate/gaugin/rest/compton_data"
	"github.com/arelate/gaugin/rest/compton_fragments"
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/direction"
	"github.com/boggydigital/compton/elements/details_toggle"
	"github.com/boggydigital/compton/elements/flex_items"
	"github.com/boggydigital/compton/elements/page"
)

func SearchNew(query map[string][]string) compton.Element {

	p := page.New("Search - gaugin").SetFavIconEmoji("🪸")

	pageStack := flex_items.New(p, direction.Column)
	p.Append(pageStack)

	appNavLinks := compton_fragments.AppNavLinks(p, compton_data.AppNavSearch)
	searchLinks := compton_fragments.SearchLinks(p, compton_data.SearchNew)

	filterSearchDetails := details_toggle.NewOpen(p, "Filter & Search")

	filterSearchDetails.Append(compton_fragments.SearchForm(p))

	pageStack.Append(
		appNavLinks,
		searchLinks,
		filterSearchDetails)

	return p
}

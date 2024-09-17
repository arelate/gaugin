package compton_pages

import (
	"github.com/arelate/gaugin/rest/compton_data"
	"github.com/arelate/gaugin/rest/compton_fragments"
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/align"
	"github.com/boggydigital/compton/consts/color"
	"github.com/boggydigital/compton/consts/direction"
	"github.com/boggydigital/compton/consts/size"
	"github.com/boggydigital/compton/elements/details_summary"
	"github.com/boggydigital/compton/elements/flex_items"
	"github.com/boggydigital/kevlar"
)

const filterSearchTitle = "Filter & Search"

func Search(query map[string][]string, ids []string, from, to int, rdx kevlar.ReadableRedux) compton.Element {

	p := compton_fragments.GauginPage(compton_data.AppNavSearch)

	pageStack := flex_items.FlexItems(p, direction.Column)
	p.Append(pageStack)

	navStack := flex_items.FlexItems(p, direction.Row).
		JustifyContent(align.Center).
		AlignItems(align.Center)

	appNavLinks := compton_fragments.
		AppNavLinks(p, compton_data.AppNavSearch)
	navStack.Append(appNavLinks)

	searchScope := compton_data.SearchScopeFromQuery(query)
	searchLinks := compton_fragments.SearchLinks(p, searchScope)
	navStack.Append(searchLinks)

	pageStack.Append(navStack)

	searchQueryDisplay := compton_fragments.SearchQueryDisplay(query, p)

	filterSearchDetails := details_summary.
		Toggle(p, filterSearchTitle, len(query) == 0).
		BackgroundColor(color.Highlight).
		SummaryMarginBlockEnd(size.Normal).
		DetailsMarginBlockEnd(size.Unset)

	filterSearchDetails.Append(compton_fragments.SearchForm(p, query, searchQueryDisplay))
	pageStack.Append(filterSearchDetails)

	if searchQueryDisplay != nil {
		pageStack.Append(searchQueryDisplay)
	}

	itemsCount := compton_fragments.ItemsCount(p, from, to, len(ids))
	pageStack.Append(itemsCount)

	if len(ids) > 0 {
		productsList := compton_fragments.ProductsList(p, ids, from, to, rdx)
		pageStack.Append(productsList)
	}

	if to < len(ids) {
		pageStack.Append(compton_fragments.ShowMoreButton(p, query, to))
	}

	pageStack.Append(compton_fragments.Footer(p))

	return p
}

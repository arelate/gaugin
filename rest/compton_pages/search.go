package compton_pages

import (
	"github.com/arelate/gaugin/rest/compton_data"
	"github.com/arelate/gaugin/rest/compton_fragments"
	"github.com/arelate/gaugin/rest/gaugin_elements/product_labels"
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/color"
	"github.com/boggydigital/compton/consts/size"
	"github.com/boggydigital/compton/elements/details_summary"
	"github.com/boggydigital/compton/elements/els"
	"github.com/boggydigital/compton/elements/flex_items"
	"github.com/boggydigital/kevlar"
)

const filterSearchTitle = "Filter & Search"

func Search(query map[string][]string, ids []string, from, to int, rdx kevlar.ReadableRedux) compton.Element {

	current := compton_data.AppNavSearch
	p, pageStack := compton_fragments.AppPage(current)
	p.AppendStyle(product_labels.StyleProductLabels)

	/* Nav stack = App navigation + Search shortcuts */

	appNavLinks := compton_fragments.AppNavLinks(p, current)

	searchScope := compton_data.SearchScopeFromQuery(query)
	searchLinks := compton_fragments.SearchLinks(p, searchScope)

	pageStack.Append(flex_items.Center(p, appNavLinks, searchLinks))

	/* Filter & Search details */

	filterSearchHeading := els.HeadingText(filterSearchTitle, detailsSummaryHeadingLevel)
	filterSearchDetails := details_summary.
		Toggle(p, filterSearchHeading, len(query) == 0).
		BackgroundColor(color.Highlight).
		SummaryMarginBlockEnd(size.Normal).
		DetailsMarginBlockEnd(size.Unset)

	searchQueryDisplay := compton_fragments.SearchQueryDisplay(query, p)

	filterSearchDetails.Append(compton_fragments.SearchForm(p, query, searchQueryDisplay))
	pageStack.Append(filterSearchDetails)

	if searchQueryDisplay != nil {
		pageStack.Append(searchQueryDisplay)
	}

	/* Search results items count */

	itemsCount := compton_fragments.ItemsCount(p, from, to, len(ids))
	pageStack.Append(itemsCount)

	/* Search results product cards */

	if len(ids) > 0 {
		productsList := compton_fragments.ProductsList(p, ids, from, to, rdx)
		pageStack.Append(productsList)
	}

	/* Show more... button */

	if to < len(ids) {
		pageStack.Append(compton_fragments.ShowMoreButton(p, query, to))
	}

	/* Standard app footer */

	pageStack.Append(compton_fragments.Footer(p))

	return p
}

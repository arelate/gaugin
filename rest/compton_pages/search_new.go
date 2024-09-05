package compton_pages

import (
	"github.com/arelate/gaugin/paths"
	"github.com/arelate/gaugin/rest/compton_data"
	"github.com/arelate/gaugin/rest/compton_fragments"
	"github.com/arelate/gaugin/rest/gaugin_elements/product_card"
	"github.com/arelate/gaugin/rest/gaugin_styles"
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/direction"
	"github.com/boggydigital/compton/consts/size"
	"github.com/boggydigital/compton/elements/details_toggle"
	"github.com/boggydigital/compton/elements/els"
	"github.com/boggydigital/compton/elements/flex_items"
	"github.com/boggydigital/compton/elements/grid_items"
	"github.com/boggydigital/compton/elements/page"
	"github.com/boggydigital/kevlar"
)

func SearchNew(query map[string][]string, ids []string, from, to int, rdx kevlar.ReadableRedux) compton.Element {

	p := page.New("Search - gaugin").
		SetFavIconEmoji("🪸").
		SetCustomStyles(gaugin_styles.GauginStyle)

	pageStack := flex_items.New(p, direction.Column)
	p.Append(pageStack)

	appNavLinks := compton_fragments.AppNavLinks(p, compton_data.AppNavSearch)
	pageStack.Append(appNavLinks)

	searchScope := compton_data.SearchScopeFromQuery(query)
	searchLinks := compton_fragments.SearchLinks(p, searchScope)
	pageStack.Append(searchLinks)

	searchQueryDisplay := compton_fragments.SearchQueryDisplay(query, p)

	var detailsToggle func(compton.Registrar, string) *details_toggle.Details
	if len(ids) > 0 {
		detailsToggle = details_toggle.NewClosed
	} else {
		detailsToggle = details_toggle.NewOpen
	}
	filterSearchDetails := detailsToggle(p, "Filter & Search")
	filterSearchDetails.Append(compton_fragments.SearchForm(p, query, searchQueryDisplay))
	pageStack.Append(filterSearchDetails)

	if searchQueryDisplay != nil {
		pageStack.Append(searchQueryDisplay)
	}

	itemsCount := compton_fragments.ItemsCount(p, from, to, len(ids))
	pageStack.Append(itemsCount)

	productCards := grid_items.New(p).
		SetColumnGap(size.Normal).
		SetRowGap(size.Normal)

	for ii := from; ii < to; ii++ {
		id := ids[ii]
		productLink := els.NewA(paths.ProductId(id))

		productCard := product_card.New(p, id, rdx)
		productLink.Append(productCard)
		productCards.Append(productLink)
	}

	pageStack.Append(productCards)

	if to < len(ids) {
		pageStack.Append(compton_fragments.ShowMoreButton(p, query, to))
	}

	return p
}

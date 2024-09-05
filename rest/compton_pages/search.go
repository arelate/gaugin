package compton_pages

import (
	"github.com/arelate/gaugin/rest/compton_data"
	"github.com/arelate/gaugin/rest/compton_fragments"
	"github.com/arelate/gaugin/rest/gaugin_styles"
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/direction"
	"github.com/boggydigital/compton/elements/details_toggle"
	"github.com/boggydigital/compton/elements/flex_items"
	"github.com/boggydigital/compton/elements/page"
	"github.com/boggydigital/kevlar"
)

func Search(query map[string][]string, ids []string, from, to int, rdx kevlar.ReadableRedux) compton.Element {

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

	filterSearchDetails := details_toggle.NewToggle(p, "Filter & Search", len(ids) == 0)
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

package compton_pages

import (
	"github.com/arelate/gaugin/rest/compton_data"
	"github.com/arelate/gaugin/rest/compton_fragments"
	"github.com/arelate/gaugin/rest/gaugin_styles"
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/elements/details_toggle"
	"github.com/boggydigital/compton/elements/flex_items"
	"github.com/boggydigital/compton/elements/nav_links"
	"github.com/boggydigital/compton/elements/page"
	"github.com/boggydigital/kevlar"
)

func Updates(sections []string, updates map[string][]string, sectionTitles map[string]string, updateTotals map[string]int, updated string, rdx kevlar.ReadableRedux) compton.Element {

	p := page.Page("Updates - gaugin").
		SetFavIconEmoji("🪸").
		SetCustomStyles(gaugin_styles.GauginStyle)

	pageStack := flex_items.FlexItemsColumn(p)
	p.Append(pageStack)

	appNavLinks := compton_fragments.AppNavLinks(p, compton_data.AppNavUpdates)
	pageStack.Append(appNavLinks)

	order := make([]string, 0, len(sections))
	sectionLinks := make(map[string]string)
	for _, section := range sections {
		st := sectionTitles[section]
		sectionLinks[st] = "#" + st
		order = append(order, st)
	}

	sectionTargets := nav_links.TextLinks(sectionLinks, "", order...)

	sectionNav := nav_links.NavLinksTargets(p, sectionTargets...)
	pageStack.Append(sectionNav)

	var showAll compton.Element
	if hasMoreItems(sections, updates, updateTotals) {
		showAll = compton_fragments.ShowAllButton(p)
		pageStack.Append(showAll)
	}

	for _, section := range sections {

		sectionDetailsToggle := details_toggle.Open(p, sectionTitles[section])
		pageStack.Append(sectionDetailsToggle)

		sectionStack := flex_items.FlexItemsColumn(p)
		sectionDetailsToggle.Append(sectionStack)

		ids := updates[section]

		itemsCount := compton_fragments.ItemsCount(p, 0, len(ids), updateTotals[section])
		sectionStack.Append(itemsCount)

		productsList := compton_fragments.ProductsList(p, ids, 0, len(ids), rdx)
		sectionStack.Append(productsList)
	}

	if showAll != nil {
		pageStack.Append(showAll)
	}

	pageStack.Append(compton_fragments.Updated(p, updated))

	pageStack.Append(compton_fragments.Footer(p))

	return p
}

func hasMoreItems(sections []string, updates map[string][]string, updateTotals map[string]int) bool {
	for _, section := range sections {
		if len(updates[section]) < updateTotals[section] {
			return true
		}
	}
	return false
}

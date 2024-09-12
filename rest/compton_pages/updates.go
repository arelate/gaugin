package compton_pages

import (
	"github.com/arelate/gaugin/rest/compton_data"
	"github.com/arelate/gaugin/rest/compton_fragments"
	"github.com/arelate/gaugin/rest/gaugin_styles"
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/align"
	"github.com/boggydigital/compton/consts/color"
	"github.com/boggydigital/compton/consts/direction"
	"github.com/boggydigital/compton/consts/size"
	"github.com/boggydigital/compton/elements/details_summary"
	"github.com/boggydigital/compton/elements/flex_items"
	"github.com/boggydigital/compton/elements/nav_links"
	"github.com/boggydigital/compton/page"
	"github.com/boggydigital/kevlar"
)

func Updates(sections []string, updates map[string][]string, sectionTitles map[string]string, updateTotals map[string]int, updated string, rdx kevlar.ReadableRedux) compton.Element {

	p := page.Page("Updates - gaugin").
		SetFavIconEmoji("🪸").
		SetCustomStyles(gaugin_styles.GauginStyle)

	pageStack := flex_items.FlexItems(p, direction.Column)
	p.Append(pageStack)

	navStack := flex_items.FlexItems(p, direction.Row).JustifyContent(align.Center).AlignItems(align.Center)

	appNavLinks := compton_fragments.AppNavLinks(p, compton_data.AppNavUpdates)
	navStack.Append(appNavLinks)

	order := make([]string, 0, len(sections))
	sectionLinks := make(map[string]string)
	for _, section := range sections {
		st := sectionTitles[section]
		sectionLinks[st] = "#" + st
		order = append(order, st)
	}

	sectionTargets := nav_links.TextLinks(sectionLinks, "", order...)

	sectionNav := nav_links.NavLinksTargets(p, sectionTargets...)
	navStack.Append(sectionNav)

	pageStack.Append(navStack)

	var showAll compton.Element
	if hasMoreItems(sections, updates, updateTotals) {
		showAll = compton_fragments.ShowAllButton(p)
		pageStack.Append(showAll)
	}

	for _, section := range sections {

		sectionDetailsToggle := details_summary.
			Open(p, sectionTitles[section]).
			BackgroundColor(color.Highlight).
			SummaryMarginBlockEnd(size.Normal).
			DetailsMarginBlockEnd(size.Large)
		pageStack.Append(sectionDetailsToggle)

		sectionStack := flex_items.FlexItems(p, direction.Column)
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

package compton_pages

import (
	"github.com/arelate/gaugin/rest/compton_data"
	"github.com/arelate/gaugin/rest/compton_fragments"
	"github.com/arelate/gaugin/rest/gaugin_elements/product_card"
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/direction"
	"github.com/boggydigital/compton/elements/details_toggle"
	"github.com/boggydigital/compton/elements/els"
	"github.com/boggydigital/compton/elements/flex_items"
	"github.com/boggydigital/compton/elements/grid_items"
	"github.com/boggydigital/compton/elements/page"
)

func SearchNew(query map[string][]string) compton.Element {

	p := page.New("Search - gaugin").SetFavIconEmoji("🪸")

	pageStack := flex_items.New(p, direction.Column)
	p.Append(pageStack)

	appNavLinks := compton_fragments.AppNavLinks(p, compton_data.AppNavSearch)
	searchLinks := compton_fragments.SearchLinks(p, compton_data.SearchNew)

	filterSearchDetails := details_toggle.NewClosed(p, "Filter & Search")

	filterSearchDetails.Append(compton_fragments.SearchForm(p))

	productCards := grid_items.New(p)

	testProductCard := product_card.New(p, "Anger Foot").
		SetDevelopers("Free Lives").
		SetPublishers("Devolver Digital").
		SetPoster(
			"21=30=I/wABCBTYZmCkFUkEJhwI4ABDSZKStElia46tSJIutYlUgiFBXpLa2LIlKdKcSxVN2uooMABFSRdHloyEcQ7GSAMDBICI8pKkJxAh2rpEFCeAFS5tgTsgwwodcFVgPnlS4slJSXPUBVghj45QWytKpLvkhufIOSRWHAgQaaS8KhZBSpKHydaTSAsBOGhwo+5TSew+zEkX44mtNm0CDNQAAgRMeeZspUsHLgPCwU8OOAQw4UOGS5E+fDA3B9yEJ0lipJuDVDEJBKop32ARzgEFhElKmJjDooTiFQ1izJlDp40SEyYElEiSJIAJdUpKAG/gIMOE4ElM5A5AYjuLPSXCr/9lAUJ0hgwHSggIgKS9gHDgwwc44KBxY+vrS7QnEQB+eBOa1WdfBg0EoB9zJRyQQ3wsNLCCA4zdtx4JSLSRjmYIuKNEEiwYAEAAEYKQAQK5LbcCOOqEsw4UVQnUQIR8TaAZCyxAQc8eV7jDom8NUQCCBiZURqB0JkBhQjiC7KhWAPUBeQOBOgWAwAE4iNGOOlCwEBaTIFBQwpMTRDnBBF3ocZwJ6Rh5AAUfeAlOAw3QiKE8x+kExRNK0AcCAtKtQGNrJtBhQpQBPAEFAA2MmCCNuiFnwhMG1tlnFTEY6KBvjkKBlBJKRNHCoMuRAAAJT1RqoE4AOucociyNmlqUBqwVtRYChALQETa2lqDCWgB4mJNOAwUEADs",
			"https://gaugin.frmnt.io/image?id=0d9684e197ff3a8d34bddab41e2ef8c9f6d1050242b44b56dfab11ff69b670bb")

	productLink := els.NewA("https://gaugin.frmnt.io/product?id=1983918143")
	productLink.Append(testProductCard)

	productCards.Append(productLink)

	pageStack.Append(
		appNavLinks,
		searchLinks,
		filterSearchDetails,
		productCards)

	return p
}

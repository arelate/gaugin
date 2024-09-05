package compton_fragments

import (
	"github.com/arelate/gaugin/paths"
	"github.com/arelate/gaugin/rest/gaugin_elements/product_card"
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/elements/els"
	"github.com/boggydigital/compton/elements/grid_items"
	"github.com/boggydigital/kevlar"
)

func ProductsList(r compton.Registrar, ids []string, from, to int, rdx kevlar.ReadableRedux) compton.Element {
	productCards := grid_items.GridItems(r)

	for ii := from; ii < to; ii++ {
		id := ids[ii]
		productLink := els.A(paths.ProductId(id))

		productCard := product_card.ProductCard(r, id, rdx)
		productLink.Append(productCard)
		productCards.Append(productLink)
	}

	return productCards
}

package compton_fragments

import (
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/color"
	"github.com/boggydigital/compton/consts/font_weight"
	"github.com/boggydigital/compton/consts/size"
	"github.com/boggydigital/compton/elements/flex_items"
	"github.com/boggydigital/compton/elements/fspan"
)

func Updated(r compton.Registrar, updated string) compton.Element {
	span := fspan.Text(r, "").
		FontSize(size.Small)
	updatedTitle := fspan.Text(r, "Updated: ").
		ForegroundColor(color.Gray)
	updatedValue := fspan.Text(r, updated).
		FontWeight(font_weight.Bolder)
	span.Append(updatedTitle, updatedValue)

	return flex_items.Center(r, span)
}

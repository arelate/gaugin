package compton_fragments

import (
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/align"
	"github.com/boggydigital/compton/consts/color"
	"github.com/boggydigital/compton/consts/direction"
	"github.com/boggydigital/compton/consts/size"
	"github.com/boggydigital/compton/consts/weight"
	"github.com/boggydigital/compton/elements/flex_items"
	"github.com/boggydigital/compton/elements/fspan"
)

func Updated(r compton.Registrar, updated string) compton.Element {
	row := flex_items.FlexItems(r, direction.Row).JustifyContent(align.Center)

	span := fspan.Text(r, "").
		FontSize(size.XSmall)
	updatedTitle := fspan.Text(r, "Updated: ").
		ForegroundColor(color.Subtle)
	updatedValue := fspan.Text(r, updated).
		FontWeight(weight.Bolder)
	span.Append(updatedTitle, updatedValue)
	row.Append(span)

	return row
}

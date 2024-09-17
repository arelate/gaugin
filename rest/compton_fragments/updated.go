package compton_fragments

import (
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/color"
	"github.com/boggydigital/compton/consts/size"
	"github.com/boggydigital/compton/consts/weight"
	"github.com/boggydigital/compton/elements/fspan"
	"github.com/boggydigital/compton/elements/recipes"
)

func Updated(r compton.Registrar, updated string) compton.Element {
	span := fspan.Text(r, "").
		FontSize(size.XSmall)
	updatedTitle := fspan.Text(r, "Updated: ").
		ForegroundColor(color.Subtle)
	updatedValue := fspan.Text(r, updated).
		FontWeight(weight.Bolder)
	span.Append(updatedTitle, updatedValue)

	return recipes.JustifyCenter(r, span)
}

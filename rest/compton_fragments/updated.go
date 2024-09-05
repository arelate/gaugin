package compton_fragments

import (
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/alignment"
	"github.com/boggydigital/compton/consts/direction"
	"github.com/boggydigital/compton/elements/els"
	"github.com/boggydigital/compton/elements/flex_items"
)

func Updated(r compton.Registrar, updated string) compton.Element {
	row := flex_items.New(r, direction.Row).JustifyContent(alignment.Center)

	span := els.NewSpan()
	span.SetClass("fs-xs")
	updatedTitle := els.NewSpanText("Updated: ")
	updatedTitle.SetClass("fg-subtle")
	updatedValue := els.NewSpanText(updated)
	span.Append(updatedTitle, updatedValue)
	row.Append(span)

	return row
}

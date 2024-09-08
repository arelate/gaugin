package compton_fragments

import (
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/align"
	"github.com/boggydigital/compton/consts/direction"
	"github.com/boggydigital/compton/elements/els"
	"github.com/boggydigital/compton/elements/flex_items"
)

func Updated(r compton.Registrar, updated string) compton.Element {
	row := flex_items.FlexItems(r, direction.Row).JustifyContent(align.Center)

	span := els.Span()
	span.AddClass("fs-xs")
	updatedTitle := els.SpanText("Updated: ")
	updatedTitle.AddClass("fg-subtle")
	updatedValue := els.SpanText(updated)
	span.Append(updatedTitle, updatedValue)
	row.Append(span)

	return row
}

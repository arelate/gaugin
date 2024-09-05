package compton_fragments

import (
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/alignment"
	"github.com/boggydigital/compton/elements/els"
	"github.com/boggydigital/compton/elements/flex_items"
)

func Updated(r compton.Registrar, updated string) compton.Element {
	row := flex_items.FlexItemsRow(r).JustifyContent(alignment.Center)

	span := els.Span()
	span.SetClass("fs-xs")
	updatedTitle := els.SpanText("Updated: ")
	updatedTitle.SetClass("fg-subtle")
	updatedValue := els.SpanText(updated)
	span.Append(updatedTitle, updatedValue)
	row.Append(span)

	return row
}

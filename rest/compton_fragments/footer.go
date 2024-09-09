package compton_fragments

import (
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/direction"
	"github.com/boggydigital/compton/consts/size"
	"github.com/boggydigital/compton/elements/els"
	"github.com/boggydigital/compton/elements/flex_items"
	"github.com/boggydigital/compton/elements/section"
)

func Footer(r compton.Registrar) compton.Element {
	sh := section.Section(r)
	sh.AddClass("footer", "fs-xs")

	row := flex_items.FlexItems(r, direction.Row).ColumnGap(size.XSmall)
	sh.Append(row)

	row.Append(
		els.SpanText("👋"), els.SpanText("from"),
		els.AText("Arles", "https://github.com/arelate"),
		els.SpanText("🇫🇷"))

	return sh
}

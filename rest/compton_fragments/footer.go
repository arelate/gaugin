package compton_fragments

import (
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/size"
	"github.com/boggydigital/compton/elements/els"
	"github.com/boggydigital/compton/elements/flex_items"
	"github.com/boggydigital/compton/elements/section_highlight"
)

func Footer(r compton.Registrar) compton.Element {
	sh := section_highlight.SectionHighlight(r)
	sh.SetClass("footer", "fs-xs")

	row := flex_items.FlexItemsRow(r).SetColumnGap(size.XSmall)
	sh.Append(row)

	row.Append(
		els.SpanText("👋"), els.SpanText("from"),
		els.AText("Arles", "https://github.com/arelate"),
		els.SpanText("🇫🇷"))

	return sh
}

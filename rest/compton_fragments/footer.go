package compton_fragments

import (
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/color"
	"github.com/boggydigital/compton/consts/direction"
	"github.com/boggydigital/compton/consts/size"
	"github.com/boggydigital/compton/consts/weight"
	"github.com/boggydigital/compton/elements/els"
	"github.com/boggydigital/compton/elements/flex_items"
	"github.com/boggydigital/compton/elements/fspan"
	"github.com/boggydigital/compton/elements/section"
)

func Footer(r compton.Registrar) compton.Element {
	sh := section.Section(r).
		BackgroundColor(color.Highlight).
		FontSize(size.Small)

	row := flex_items.FlexItems(r, direction.Row).ColumnGap(size.XSmall)
	sh.Append(row)

	link := els.A("https://github.com/arelate")
	link.Append(fspan.Text(r, "Arles").FontWeight(weight.Bolder))

	row.Append(
		els.SpanText("👋"), els.SpanText("from"),
		link,
		els.SpanText("🇫🇷"))

	return sh
}

package compton_fragments

import (
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/font_weight"
	"github.com/boggydigital/compton/consts/size"
	"github.com/boggydigital/compton/elements/els"
	"github.com/boggydigital/compton/elements/flex_items"
	"github.com/boggydigital/compton/elements/fspan"
)

func Footer(r compton.Registrar) compton.Element {

	link := els.A("https://github.com/arelate")
	link.Append(fspan.Text(r, "Arles").FontWeight(font_weight.Bolder).FontSize(size.Small))

	row := flex_items.Center(r,
		fspan.Text(r, "👋").FontSize(size.Small),
		fspan.Text(r, "from").FontSize(size.Small),
		link,
		fspan.Text(r, "🇫🇷").FontSize(size.Small))

	row.ColumnGap(size.XSmall)

	return row
}

package compton_fragments

import (
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/direction"
	"github.com/boggydigital/compton/consts/size"
	"github.com/boggydigital/compton/elements/els"
	"github.com/boggydigital/compton/elements/flex_items"
	"github.com/boggydigital/compton/elements/section_highlight"
)

func Footer(r compton.Registrar) compton.Element {
	sh := section_highlight.New(r)
	sh.SetClass("footer", "fs-xs")

	row := flex_items.New(r, direction.Row).SetColumnGap(size.XSmall)
	sh.Append(row)

	hello := els.NewSpanText("👋")
	from := els.NewSpanText("from")
	arlesFrance := els.NewAText("Arles 🇫🇷", "https://github.com/arelate")
	row.Append(hello, from, arlesFrance)

	return sh
}

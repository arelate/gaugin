package compton_fragments

import (
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/size"
	"github.com/boggydigital/compton/consts/weight"
	"github.com/boggydigital/compton/elements/fspan"
)

func DetailsSummaryTitle(r compton.Registrar, title string) compton.Element {
	fs := fspan.Text(r, title).
		FontWeight(weight.Bolder).
		FontSize(size.Large)
	return fs
}

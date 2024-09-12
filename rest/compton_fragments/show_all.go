package compton_fragments

import (
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/align"
	"github.com/boggydigital/compton/consts/direction"
	"github.com/boggydigital/compton/consts/input_types"
	"github.com/boggydigital/compton/elements/els"
	"github.com/boggydigital/compton/elements/flex_items"
	"github.com/boggydigital/compton/elements/inputs"
)

func ShowAllButton(r compton.Registrar) compton.Element {

	row := flex_items.FlexItems(r, direction.Row).JustifyContent(align.Center)

	showAllLink := els.A("?show-all=true")
	row.Append(showAllLink)

	button := inputs.InputValue(r, input_types.Submit, "Show All...")
	showAllLink.Append(button)

	return row

}

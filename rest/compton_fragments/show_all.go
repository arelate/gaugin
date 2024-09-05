package compton_fragments

import (
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/alignment"
	"github.com/boggydigital/compton/consts/input_types"
	"github.com/boggydigital/compton/elements/els"
	"github.com/boggydigital/compton/elements/flex_items"
)

func ShowAllButton(r compton.Registrar) compton.Element {

	row := flex_items.FlexItemsRow(r).JustifyContent(alignment.Center)

	showAllLink := els.A("?show-all=true")
	showAllLink.SetClass("updates-show-all")
	row.Append(showAllLink)

	button := els.InputValue(input_types.Submit, "Show All...")
	showAllLink.Append(button)

	return row

}

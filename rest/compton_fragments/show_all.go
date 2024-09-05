package compton_fragments

import (
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/alignment"
	"github.com/boggydigital/compton/consts/direction"
	"github.com/boggydigital/compton/consts/input_types"
	"github.com/boggydigital/compton/elements/els"
	"github.com/boggydigital/compton/elements/flex_items"
)

func ShowAllButton(r compton.Registrar) compton.Element {

	row := flex_items.New(r, direction.Row).JustifyContent(alignment.Center)

	showAllLink := els.NewA("?show-all=true")
	showAllLink.SetClass("search-show-more")
	row.Append(showAllLink)

	button := els.NewInputValue(input_types.Submit, "Show All...")
	showAllLink.Append(button)

	return row

}

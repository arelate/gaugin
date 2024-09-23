package compton_fragments

import (
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/input_types"
	"github.com/boggydigital/compton/elements/els"
	"github.com/boggydigital/compton/elements/flex_items"
	"github.com/boggydigital/compton/elements/inputs"
)

func ShowAllButton(r compton.Registrar) compton.Element {

	showAllLink := els.A("?show-all=true")

	button := inputs.InputValue(r, input_types.Submit, "Show all...")
	showAllLink.Append(button)

	return flex_items.Center(r, showAllLink)

}

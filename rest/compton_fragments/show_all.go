package compton_fragments

import (
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/input_types"
	"github.com/boggydigital/compton/elements/els"
	"github.com/boggydigital/compton/elements/inputs"
	"github.com/boggydigital/compton/elements/recipes"
)

func ShowAllButton(r compton.Registrar) compton.Element {

	showAllLink := els.A("?show-all=true")

	button := inputs.InputValue(r, input_types.Submit, "Show All...")
	showAllLink.Append(button)

	return recipes.Center(r, showAllLink)

}

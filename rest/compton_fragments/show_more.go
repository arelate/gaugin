package compton_fragments

import (
	"github.com/arelate/gaugin/rest/compton_data"
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/input_types"
	"github.com/boggydigital/compton/elements/els"
	"github.com/boggydigital/compton/elements/inputs"
	"github.com/boggydigital/compton/elements/recipes"
	"strconv"
)

func ShowMoreButton(r compton.Registrar, query map[string][]string, from int) compton.Element {

	query["from"] = []string{strconv.Itoa(from)}
	enq := compton_data.EncodeQuery(query)

	showMoreLink := els.A("/search?" + enq)
	button := inputs.InputValue(r, input_types.Submit, "Show more...")
	showMoreLink.Append(button)

	return recipes.Center(r, showMoreLink)
}

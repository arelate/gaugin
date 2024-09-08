package compton_fragments

import (
	"github.com/arelate/gaugin/rest/compton_data"
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/align"
	"github.com/boggydigital/compton/consts/direction"
	"github.com/boggydigital/compton/consts/input_types"
	"github.com/boggydigital/compton/elements/els"
	"github.com/boggydigital/compton/elements/flex_items"
	"strconv"
)

func ShowMoreButton(r compton.Registrar, query map[string][]string, from int) compton.Element {

	query["from"] = []string{strconv.Itoa(from)}
	enq := compton_data.EncodeQuery(query)

	row := flex_items.FlexItems(r, direction.Row).JustifyContent(align.Center)

	showMoreLink := els.A("/search?" + enq)
	showMoreLink.AddClass("search-show-more")
	row.Append(showMoreLink)

	button := els.InputValue(input_types.Submit, "More...")
	showMoreLink.Append(button)

	return row

}

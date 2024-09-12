package compton_fragments

import (
	"github.com/arelate/gaugin/rest/compton_data"
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/align"
	"github.com/boggydigital/compton/consts/direction"
	"github.com/boggydigital/compton/consts/input_types"
	"github.com/boggydigital/compton/elements/els"
	"github.com/boggydigital/compton/elements/flex_items"
	"github.com/boggydigital/compton/elements/grid_items"
	"github.com/boggydigital/compton/elements/inputs"
	"github.com/boggydigital/compton/elements/title_values"
	"slices"
	"strings"
)

func SearchForm(r compton.Registrar, query map[string][]string, searchQueryDisplay compton.Element) compton.Element {

	form := els.Form("/search", "GET")
	formStack := flex_items.FlexItems(r, direction.Column)
	form.Append(formStack)

	if searchQueryDisplay != nil {
		formStack.Append(searchQueryDisplay)
	}

	submitRow := flex_items.FlexItems(r, direction.Row).JustifyContent(align.Center)
	submit := inputs.InputValue(r, input_types.Submit, "Submit Query")
	submitRow.Append(submit)
	formStack.Append(submitRow)

	inputsGrid := grid_items.GridItems(r).JustifyContent(align.Center)
	formStack.Append(inputsGrid)

	searchInputs(r, query, inputsGrid)

	// duplicating Submit button after inputs at the end
	formStack.Append(submitRow)

	return form
}

func searchInputs(r compton.Registrar, query map[string][]string, container compton.Element) {
	for _, property := range compton_data.SearchProperties {
		title := compton_data.PropertyTitles[property]
		value := strings.Join(query[property], ", ")
		titleInput := title_values.SearchValue(r, title, property, value)
		if slices.Contains(compton_data.DigestProperties, property) {
			// set datalist for that property
		}
		container.Append(titleInput)
	}
}

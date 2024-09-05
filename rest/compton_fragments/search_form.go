package compton_fragments

import (
	"github.com/arelate/gaugin/rest/compton_data"
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/alignment"
	"github.com/boggydigital/compton/consts/direction"
	"github.com/boggydigital/compton/consts/input_types"
	"github.com/boggydigital/compton/elements/els"
	"github.com/boggydigital/compton/elements/flex_items"
	"github.com/boggydigital/compton/elements/grid_items"
	"github.com/boggydigital/compton/elements/title_values"
	"slices"
)

func SearchForm(searchQuery compton.Element, r compton.Registrar) compton.Element {

	form := els.NewForm("/search", "GET")
	formStack := flex_items.New(r, direction.Column)
	form.Append(formStack)

	if searchQuery != nil {
		formStack.Append(searchQuery)
	}

	submitRow := flex_items.New(r, direction.Row).
		JustifyContent(alignment.Center)
	submit := els.NewInputValue(input_types.Submit, "Submit Query")
	submitRow.Append(submit)
	formStack.Append(submitRow)

	inputsGrid := grid_items.New(r)
	formStack.Append(inputsGrid)

	searchInputs(r, inputsGrid)

	// duplicating Submit button after inputs at the end
	formStack.Append(submitRow)

	return form
}

func searchInputs(r compton.Registrar, container compton.Element) {
	for _, property := range compton_data.SearchProperties {
		title := compton_data.PropertyTitles[property]
		titleInput := title_values.NewSearchValue(r, title, property, "")
		if slices.Contains(compton_data.DigestProperties, property) {
			// set datalist for that property
		}
		container.Append(titleInput)
	}
}

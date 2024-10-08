package compton_fragments

import (
	"fmt"
	"github.com/arelate/gaugin/rest/compton_data"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/align"
	"github.com/boggydigital/compton/consts/direction"
	"github.com/boggydigital/compton/consts/input_types"
	"github.com/boggydigital/compton/elements/els"
	"github.com/boggydigital/compton/elements/flex_items"
	"github.com/boggydigital/compton/elements/grid_items"
	"github.com/boggydigital/compton/elements/inputs"
	"github.com/boggydigital/compton/elements/title_values"
	"golang.org/x/exp/maps"
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

var binDatalist = map[string]string{
	"true":  "Yes",
	"false": "No",
}

var typesDigest = []vangogh_local_data.ProductType{
	vangogh_local_data.AccountProducts,
	vangogh_local_data.CatalogProducts,
}

func stringerDatalist[T fmt.Stringer](items []T) map[string]string {
	dl := make(map[string]string)
	for _, item := range items {
		str := item.String()
		dl[str] = compton_data.PropertyTitles[str]
	}
	return dl
}

func typesDatalist() map[string]string {
	return stringerDatalist(typesDigest)
}

func operatingSystemsDatalist() map[string]string {
	return stringerDatalist([]vangogh_local_data.OperatingSystem{
		vangogh_local_data.Windows,
		vangogh_local_data.MacOS,
		vangogh_local_data.Linux})
}

var sortProperties = []string{
	vangogh_local_data.GlobalReleaseDateProperty,
	vangogh_local_data.GOGReleaseDateProperty,
	vangogh_local_data.GOGOrderDateProperty,
	vangogh_local_data.TitleProperty,
	vangogh_local_data.RatingProperty,
	vangogh_local_data.DiscountPercentageProperty,
	vangogh_local_data.HLTBHoursToCompleteMainProperty,
	vangogh_local_data.HLTBHoursToCompletePlusProperty,
	vangogh_local_data.HLTBHoursToComplete100Property}

func propertiesDatalist(properties []string) map[string]string {
	dl := make(map[string]string)
	for _, p := range properties {
		dl[p] = compton_data.PropertyTitles[p]
	}
	return dl
}

func sortDatalist() map[string]string {
	return propertiesDatalist(sortProperties)
}

func productTypesDatalist() map[string]string {
	return propertiesDatalist([]string{"GAME", "PACK", "DLC"})
}

func steamDeckDatalist() map[string]string {
	return propertiesDatalist([]string{"Verified", "Playable", "Unsupported", "Unknown"})
}

func languagesDatalist() map[string]string {
	dl := make(map[string]string)
	for _, lc := range maps.Keys(compton_data.LanguageTitles) {
		dl[lc] = compton_data.FormatLanguage(lc)
	}
	return dl
}

func searchInputs(r compton.Registrar, query map[string][]string, container compton.Element) {
	for _, property := range compton_data.SearchProperties {
		title := compton_data.PropertyTitles[property]
		value := strings.Join(query[property], ", ")
		titleInput := title_values.SearchValue(r, title, property, value)

		var datalist map[string]string
		var listId string

		if slices.Contains(compton_data.BinaryDigestProperties, property) {
			datalist = binDatalist
			listId = "bin-list"
		} else if slices.Contains(compton_data.DigestProperties, property) {
			switch property {
			case vangogh_local_data.TypesProperty:
				datalist = typesDatalist()
			case vangogh_local_data.OperatingSystemsProperty:
				datalist = operatingSystemsDatalist()
			case vangogh_local_data.SortProperty:
				datalist = sortDatalist()
			case vangogh_local_data.ProductTypeProperty:
				datalist = productTypesDatalist()
			case vangogh_local_data.SteamDeckAppCompatibilityCategoryProperty:
				datalist = steamDeckDatalist()
			case vangogh_local_data.LanguageCodeProperty:
				datalist = languagesDatalist()
			}
		}

		if len(datalist) > 0 {
			titleInput.SetDatalist(datalist, listId)
		}

		container.Append(titleInput)
	}
}

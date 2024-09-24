package compton_pages

import (
	"github.com/arelate/gaugin/rest/compton_data"
	"github.com/arelate/gaugin/rest/compton_fragments"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/align"
	"github.com/boggydigital/compton/consts/color"
	"github.com/boggydigital/compton/consts/direction"
	"github.com/boggydigital/compton/consts/input_types"
	"github.com/boggydigital/compton/consts/size"
	"github.com/boggydigital/compton/elements/details_summary"
	"github.com/boggydigital/compton/elements/els"
	"github.com/boggydigital/compton/elements/flex_items"
	"github.com/boggydigital/compton/elements/inputs"
	"github.com/boggydigital/kevlar"
	"golang.org/x/exp/maps"
	"net/http"
	"slices"
)

func TagsEditor(
	id string,
	owned bool,
	tagsProperty string,
	allValues map[string]string,
	selected map[string]any,
	rdx kevlar.ReadableRedux) compton.Element {

	tagsPropertyTitle := compton_data.PropertyTitles[tagsProperty]

	p, pageStack := compton_fragments.AppPage("Edit " + tagsPropertyTitle)

	/* App navigation */

	appNavLinks := compton_fragments.AppNavLinks(p, "")

	pageStack.Append(flex_items.Center(p, appNavLinks))

	/* Product poster */

	if poster := compton_fragments.ProductPoster(p, id, rdx); poster != nil {
		pageStack.Append(poster)
	}

	/* Product title */

	productTitle, _ := rdx.GetLastVal(vangogh_local_data.TitleProperty, id)
	productHeading := els.HeadingText(productTitle, 1)
	pageStack.Append(flex_items.Center(p, productHeading))

	/* Tags Property Title */

	tagsPropertyHeading := compton_fragments.DetailsSummaryTitle(p, tagsPropertyTitle)

	dsTags := details_summary.Open(p, tagsPropertyHeading).
		BackgroundColor(color.Indigo).
		ForegroundColor(color.Background).
		MarkerColor(color.Background).
		DetailsMarginBlockEnd(size.Large)

	pageStack.Append(dsTags)

	/* Tag Values Switches */

	editTagsForm := els.Form("/local-tags/apply", http.MethodGet)
	swColumn := flex_items.FlexItems(p, direction.Column).AlignContent(align.Center)

	idInput := inputs.InputValue(p, input_types.Hidden, id)
	idInput.SetName(vangogh_local_data.IdProperty)
	swColumn.Append(idInput)

	keys := maps.Keys(allValues)
	slices.Sort(keys)

	for _, vid := range keys {
		label := allValues[vid]
		_, has := selected[vid]
		swColumn.Append(switchLabel(p, label, has))
	}

	newValueInput := inputs.Input(p, input_types.Text)
	newValueInput.SetName("new-property-value")
	newValueInput.SetPlaceholder("Add new value")
	swColumn.Append(newValueInput)

	applyButton := inputs.InputValue(p, input_types.Submit, "Apply")
	swColumn.Append(applyButton)

	editTagsForm.Append(swColumn)
	dsTags.Append(editTagsForm)

	/* Footer */

	pageStack.Append(compton_fragments.Footer(p))

	return p
}

func switchLabel(r compton.Registrar, label string, checked bool) compton.Element {
	row := flex_items.FlexItems(r, direction.Row).AlignItems(align.Center)

	switchElement := inputs.Switch(r)
	switchElement.SetId(label)
	switchElement.SetValue(label)
	switchElement.SetChecked(checked)
	switchElement.SetName("value") //using the same name for all binary properties

	labelElement := els.Label(label)
	labelElement.Append(els.Text(label))

	row.Append(switchElement, labelElement)

	return row
}

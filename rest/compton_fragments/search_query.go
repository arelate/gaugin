package compton_fragments

import (
	"github.com/arelate/gaugin/rest/compton_data"
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/align"
	"github.com/boggydigital/compton/consts/direction"
	"github.com/boggydigital/compton/consts/size"
	"github.com/boggydigital/compton/elements/els"
	"github.com/boggydigital/compton/elements/flex_items"
	"github.com/boggydigital/compton/elements/section"
	"golang.org/x/exp/maps"
	"slices"
	"strings"
)

func SearchQueryDisplay(query map[string][]string, r compton.Registrar) compton.Element {
	if len(query) == 0 {
		return nil
	}

	sh := section.Section(r)
	sh.AddClass("fs-xs")

	shStack := flex_items.FlexItems(r, direction.Row).
		//ColumnGap(size.Normal).
		RowGap(size.Small).
		JustifyContent(align.Center)
	sh.Append(shStack)

	sortedProperties := maps.Keys(query)
	slices.Sort(sortedProperties)

	for _, property := range sortedProperties {
		values := query[property]
		span := els.Span()
		propertyTitle := els.SpanText(compton_data.PropertyTitles[property] + ": ")
		propertyTitle.AddClass("fg-subtle")
		propertyValue := els.SpanText(strings.Join(values, ", "))
		propertyValue.AddClass("fw-b")
		span.Append(propertyTitle, propertyValue)
		shStack.Append(span)
	}

	clearAction := els.AText("Clear", "/search")
	clearAction.AddClass("action")
	shStack.Append(clearAction)

	return sh
}

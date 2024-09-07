package compton_fragments

import (
	"github.com/arelate/gaugin/rest/compton_data"
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/alignment"
	"github.com/boggydigital/compton/consts/size"
	"github.com/boggydigital/compton/elements/c_section"
	"github.com/boggydigital/compton/elements/els"
	"github.com/boggydigital/compton/elements/flex_items"
	"golang.org/x/exp/maps"
	"slices"
	"strings"
)

func SearchQueryDisplay(query map[string][]string, r compton.Registrar) compton.Element {
	if len(query) == 0 {
		return nil
	}

	sh := c_section.CSection(r)
	sh.SetClass("fs-xs")

	shStack := flex_items.FlexItemsRow(r).
		SetColumnGap(size.Normal).
		SetRowGap(size.Small).
		JustifyContent(alignment.Center)
	sh.Append(shStack)

	sortedProperties := maps.Keys(query)
	slices.Sort(sortedProperties)

	for _, property := range sortedProperties {
		values := query[property]
		span := els.Span()
		propertyTitle := els.SpanText(compton_data.PropertyTitles[property] + ": ")
		propertyTitle.SetClass("fg-subtle")
		propertyValue := els.SpanText(strings.Join(values, ", "))
		propertyValue.SetClass("fw-b")
		span.Append(propertyTitle, propertyValue)
		shStack.Append(span)
	}

	clearAction := els.AText("Clear", "/search")
	clearAction.SetClass("action")
	shStack.Append(clearAction)

	return sh
}

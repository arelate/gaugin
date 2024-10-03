package compton_fragments

import (
	"github.com/arelate/gaugin/rest/compton_data"
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/align"
	"github.com/boggydigital/compton/consts/color"
	"github.com/boggydigital/compton/consts/direction"
	"github.com/boggydigital/compton/consts/font_weight"
	"github.com/boggydigital/compton/consts/size"
	"github.com/boggydigital/compton/elements/els"
	"github.com/boggydigital/compton/elements/flex_items"
	"github.com/boggydigital/compton/elements/fspan"
	"github.com/boggydigital/compton/elements/section"
	"golang.org/x/exp/maps"
	"slices"
	"strings"
)

func SearchQueryDisplay(query map[string][]string, r compton.Registrar) compton.Element {
	if len(query) == 0 {
		return nil
	}

	sh := section.Section(r).
		BackgroundColor(color.Highlight).
		FontSize(size.Small).
		ColumnGap(size.Small)

	shStack := flex_items.FlexItems(r, direction.Row).
		RowGap(size.Small).
		JustifyContent(align.Center)
	sh.Append(shStack)

	sortedProperties := maps.Keys(query)
	slices.Sort(sortedProperties)

	for _, property := range sortedProperties {
		values := query[property]
		span := els.Span()
		propertyTitleLink := els.A("#" + compton_data.PropertyTitles[property])
		propertyTitleText := fspan.Text(r, compton_data.PropertyTitles[property]+": ").
			ForegroundColor(color.Gray)
		propertyTitleLink.Append(propertyTitleText)
		fmtValues := make([]string, 0, len(values))
		for _, value := range values {
			fmtVal := value
			if pt, ok := compton_data.PropertyTitles[value]; ok {
				fmtVal = pt
			}
			fmtValues = append(fmtValues, fmtVal)
		}
		propertyValue := fspan.Text(r, strings.Join(fmtValues, ", ")).
			FontWeight(font_weight.Bolder)
		span.Append(propertyTitleLink, propertyValue)
		shStack.Append(span)
	}

	clearLink := els.A("/search")
	clearText := fspan.Text(r, "Clear").
		ForegroundColor(color.Blue).FontWeight(font_weight.Bolder)
	clearLink.Append(clearText)
	shStack.Append(clearLink)

	return sh
}

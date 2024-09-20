package compton_fragments

import (
	"github.com/arelate/gaugin/rest/compton_data"
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/align"
	"github.com/boggydigital/compton/consts/color"
	"github.com/boggydigital/compton/consts/direction"
	"github.com/boggydigital/compton/consts/size"
	"github.com/boggydigital/compton/consts/weight"
	"github.com/boggydigital/compton/elements/els"
	"github.com/boggydigital/compton/elements/flex_items"
	"github.com/boggydigital/compton/elements/fspan"
	"github.com/boggydigital/compton/elements/section"
)

func ProductSectionsLinks(r compton.Registrar, sections []string) compton.Element {

	linksSection := section.Section(r).
		BackgroundColor(color.Highlight).
		FontSize(size.Small).
		FontWeight(weight.Normal)

	span := fspan.Text(r, "Scroll to...").
		FontSize(size.Normal).
		FontWeight(weight.Bolder)
	ds := els.Details().AppendSummary(span)
	ds.AddClass("product-sections-links")

	linksStack := flex_items.FlexItems(r, direction.Row).
		JustifyContent(align.Center).
		RowGap(size.Small)

	ds.Append(linksStack)
	linksSection.Append(ds)

	for _, s := range sections {
		title := compton_data.SectionTitles[s]
		link := els.A("#" + title)
		linkText := fspan.Text(r, title).
			ForegroundColor(color.LightBlue)
		link.Append(linkText)
		linksStack.Append(link)
	}

	return linksSection
}

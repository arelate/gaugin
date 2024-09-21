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
	"github.com/boggydigital/compton/elements/recipes"
	"github.com/boggydigital/compton/elements/section"
)

func ProductSectionsLinks(r compton.Registrar, sections []string) compton.Element {

	linksSection := section.Section(r).
		BackgroundColor(color.Highlight).
		FontSize(size.Small).
		FontWeight(weight.Normal)

	linksStack := flex_items.FlexItems(r, direction.Row).
		JustifyContent(align.Center).
		RowGap(size.Small)

	for _, s := range sections {
		title := compton_data.SectionTitles[s]
		link := els.A("#" + title)
		linkText := fspan.Text(r, title)
		link.Append(linkText)
		linksStack.Append(link)
	}

	linksSection.Append(linksStack)

	wrapper := recipes.Center(r, linksSection)
	wrapper.SetId("product-sections-links")

	return wrapper
}

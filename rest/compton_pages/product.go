package compton_pages

import (
	"github.com/arelate/gaugin/rest/compton_data"
	"github.com/arelate/gaugin/rest/compton_fragments"
	"github.com/arelate/gaugin/rest/gaugin_elements/product_labels"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/color"
	"github.com/boggydigital/compton/consts/input_types"
	"github.com/boggydigital/compton/consts/size"
	"github.com/boggydigital/compton/elements/details_summary"
	"github.com/boggydigital/compton/elements/els"
	"github.com/boggydigital/compton/elements/flex_items"
	"github.com/boggydigital/compton/elements/iframe_expand"
	"github.com/boggydigital/compton/elements/inputs"
	"github.com/boggydigital/compton/elements/labels"
	"github.com/boggydigital/compton/elements/popup"
	"github.com/boggydigital/kevlar"
	"slices"
)

var convertedSections = []string{
	compton_data.PropertiesSection,
	compton_data.ExternalLinksSection,
	compton_data.ScreenshotsSection,
	compton_data.ChangelogSection,
	compton_data.DescriptionSection,
	compton_data.VideosSection,
	compton_data.SteamDeckSection,
}

func Product(id string, rdx kevlar.ReadableRedux, hasSections []string) compton.Element {

	title, ok := rdx.GetLastVal(vangogh_local_data.TitleProperty, id)
	if !ok {
		return nil
	}

	p, pageStack := compton_fragments.AppPage(title)
	p.AppendStyle(product_labels.StyleProductLabels)

	/* App navigation */

	appNavLinks := compton_fragments.AppNavLinks(p, "")

	showToc := inputs.InputValue(p, input_types.Button, "Scroll to...")
	pageStack.Append(flex_items.Center(p, appNavLinks, showToc))

	/* Product details sections shortcuts */

	productSectionsLinks := compton_fragments.ProductSectionsLinks(p, hasSections)
	pageStack.Append(productSectionsLinks)

	pageStack.Append(popup.Attach(p, showToc, productSectionsLinks))

	/* Product poster */

	if poster := compton_fragments.ProductPoster(p, id, rdx); poster != nil {
		pageStack.Append(poster)
	}

	/* Product title */

	productTitle := els.HeadingText(title, 1)
	productTitle.AddClass("product-title")

	/* Product labels */

	fmtLabels := product_labels.FormatLabels(id, rdx, compton_data.LabelProperties...)
	productLabels := labels.Labels(p, fmtLabels...).FontSize(size.Small).RowGap(size.XSmall).ColumnGap(size.XSmall)
	pageStack.Append(flex_items.Center(p, productTitle, productLabels))

	/* Product details sections */

	for _, section := range hasSections {

		dsbc := color.Highlight
		dsfc := color.Foreground
		dsmc := color.Subtle
		if !slices.Contains(convertedSections, section) {
			dsbc = color.Red
			dsfc = color.Background
			dsmc = color.Black
		}

		sectionTitle := compton_data.SectionTitles[section]
		summaryHeading := compton_fragments.DetailsSummaryTitle(p, sectionTitle)
		detailsSummary := details_summary.
			Toggle(p, summaryHeading, section == compton_data.PropertiesSection).
			BackgroundColor(dsbc).
			ForegroundColor(dsfc).
			MarkerColor(dsmc).
			SummaryMarginBlockEnd(size.Large).
			DetailsMarginBlockEnd(size.Normal)
		detailsSummary.SetId(sectionTitle)

		switch section {
		case compton_data.PropertiesSection:
			if productProperties := compton_fragments.ProductProperties(p, id, rdx); productProperties != nil {
				detailsSummary.Append(productProperties)
			}
		case compton_data.ExternalLinksSection:
			if externalLinks := compton_fragments.ProductExternalLinks(p, id, rdx); externalLinks != nil {
				detailsSummary.Append(externalLinks)
			}
		default:
			ifh := iframe_expand.IframeExpandHost(p, section, "/"+section+"?id="+id)
			detailsSummary.Append(ifh)
		}
		pageStack.Append(detailsSummary)
	}

	/* Standard app footer */

	pageStack.Append(compton_fragments.Footer(p))

	return p
}

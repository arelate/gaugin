package compton_pages

import (
	"github.com/arelate/gaugin/rest/compton_data"
	"github.com/arelate/gaugin/rest/compton_fragments"
	"github.com/arelate/gaugin/rest/gaugin_elements/product_labels"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/color"
	"github.com/boggydigital/compton/consts/size"
	"github.com/boggydigital/compton/elements/details_summary"
	"github.com/boggydigital/compton/elements/els"
	"github.com/boggydigital/compton/elements/iframe_expand"
	"github.com/boggydigital/compton/elements/issa_image"
	"github.com/boggydigital/compton/elements/recipes"
	"github.com/boggydigital/issa"
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
}

func Product(id string, rdx kevlar.ReadableRedux, hasSections []string) compton.Element {

	title, ok := rdx.GetLastVal(vangogh_local_data.TitleProperty, id)
	if !ok {
		return nil
	}

	p, pageStack := compton_fragments.AppPage(title)

	/* App navigation */

	appNavLinks := compton_fragments.AppNavLinks(p, "")
	pageStack.Append(recipes.Center(p, appNavLinks))

	/* Product poster */

	if imgSrc, ok := rdx.GetLastVal(vangogh_local_data.ImageProperty, id); ok {
		var poster compton.Element
		relImgSrc := "/image?id=" + imgSrc
		if dehydSrc, sure := rdx.GetLastVal(vangogh_local_data.DehydratedImageProperty, id); sure {
			hydSrc := issa.HydrateColor(dehydSrc)
			poster = issa_image.IssaImageHydrated(p, hydSrc, relImgSrc)
		} else {
			poster = els.Img(relImgSrc)
		}
		poster.AddClass("product-poster")
		pageStack.Append(poster)
	}

	/* Product title */

	productTitle := els.HeadingText(title, 1)
	productTitle.AddClass("product-title")
	pageStack.Append(recipes.Center(p, productTitle))

	/* Product labels */

	labels := product_labels.Labels(p, id, rdx).FontSize(size.Small).RowGap(size.XSmall).ColumnGap(size.XSmall)
	labelsCenter := recipes.Center(p, labels)
	labelsCenter.AddClass("labels")
	pageStack.Append(labelsCenter)

	/* Product details sections shortcuts */

	pageStack.Append(compton_fragments.ProductSectionsLinks(p, hasSections))

	/* Product details sections */

	for _, section := range hasSections {

		dsbc := color.Highlight
		dsfc := color.Foreground
		if !slices.Contains(convertedSections, section) {
			dsbc = color.Red
			dsfc = color.Background
		}

		detailsSummary := details_summary.
			Toggle(p, compton_data.SectionTitles[section], section == compton_data.PropertiesSection).
			BackgroundColor(dsbc).
			ForegroundColor(dsfc).
			SummaryMarginBlockEnd(size.Large).
			DetailsMarginBlockEnd(size.Normal)

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

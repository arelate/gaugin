package compton_pages

import (
	"github.com/arelate/gaugin/rest/compton_data"
	"github.com/arelate/gaugin/rest/compton_fragments"
	"github.com/arelate/gaugin/rest/gaugin_elements/product_labels"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/color"
	"github.com/boggydigital/compton/elements/details_summary"
	"github.com/boggydigital/compton/elements/els"
	"github.com/boggydigital/compton/elements/iframe_expand"
	"github.com/boggydigital/compton/elements/issa_image"
	"github.com/boggydigital/compton/elements/recipes"
	"github.com/boggydigital/issa"
	"github.com/boggydigital/kevlar"
)

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
		if dehydSrc, sure := rdx.GetLastVal(vangogh_local_data.DehydratedImageProperty, id); sure {
			hydSrc := issa.HydrateColor(dehydSrc)
			poster = issa_image.IssaImageHydrated(p, hydSrc, imgSrc)
		} else {
			poster = els.Img(imgSrc)
		}
		poster.AddClass("product-poster")
		pageStack.Append(poster)
	}

	/* Product title */

	productTitle := els.HeadingText(title, 1)
	productTitle.AddClass("product-title")
	pageStack.Append(recipes.Center(p, productTitle))

	/* Product labels */

	labels := recipes.Center(p, product_labels.Labels(p, id, rdx))
	labels.AddClass("labels")
	pageStack.Append(labels)

	/* Product details sections shortcuts */

	pageStack.Append(compton_fragments.ProductSectionsLinks(p, hasSections))

	/* Product details sections */

	for _, section := range hasSections {
		detailsSummary := details_summary.
			Toggle(p, compton_data.SectionTitles[section], toggleProductSection(section)).
			BackgroundColor(color.Highlight)
		switch section {
		case compton_data.PropertiesSection:
			if productProperties := compton_fragments.ProductProperties(p, id, rdx); productProperties != nil {
				detailsSummary.Append(productProperties)
			}
		case compton_data.LinksSection:
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

func toggleProductSection(section string) bool {
	if section == compton_data.PropertiesSection {
		return true
	}
	return false
}

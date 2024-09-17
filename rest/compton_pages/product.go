package compton_pages

import (
	"github.com/arelate/gaugin/rest/compton_fragments"
	"github.com/arelate/gaugin/rest/gaugin_elements/product_labels"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/align"
	"github.com/boggydigital/compton/consts/direction"
	"github.com/boggydigital/compton/elements/els"
	"github.com/boggydigital/compton/elements/flex_items"
	"github.com/boggydigital/compton/elements/issa_image"
	"github.com/boggydigital/compton/elements/recipes"
	"github.com/boggydigital/issa"
	"github.com/boggydigital/kevlar"
)

func Product(id string, rdx kevlar.ReadableRedux) compton.Element {

	title, ok := rdx.GetLastVal(vangogh_local_data.TitleProperty, id)
	if !ok {
		return nil
	}

	p, pageStack := compton_fragments.AppPage(title)

	/* App navigation */

	navStack := flex_items.FlexItems(p, direction.Row).
		JustifyContent(align.Center).
		AlignItems(align.Center)
	pageStack.Append(navStack)

	appNavLinks := compton_fragments.AppNavLinks(p, "")
	navStack.Append(appNavLinks)

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

	pageStack.Append(
		recipes.JustifyCenter(p, els.HeadingText(title, 1)))

	/* Product labels */

	labels := recipes.JustifyCenter(p, product_labels.Labels(p, id, rdx))
	labels.AddClass("labels")
	pageStack.Append(labels)

	/* Product details sections shortcuts */

	/* Product details sections */

	/* Standard app footer */

	pageStack.Append(compton_fragments.Footer(p))

	return p
}

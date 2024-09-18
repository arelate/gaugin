package compton_fragments

import (
	"github.com/arelate/gaugin/rest/compton_data"
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/align"
	"github.com/boggydigital/compton/elements/grid_items"
	"github.com/boggydigital/compton/elements/title_values"
	"github.com/boggydigital/kevlar"
	"strings"
)

func ProductExternalLinks(r compton.Registrar, id string, rdx kevlar.ReadableRedux) compton.Element {

	grid := grid_items.GridItems(r).JustifyContent(align.Center)

	for _, linkProperty := range compton_data.ProductExternalLinksProperties {
		if links, ok := rdx.GetAllValues(linkProperty, id); ok && len(links) > 0 {
			if extLinks := externalLinks(r, linkProperty, links); extLinks != nil {
				grid.Append(extLinks)
			}
		}
	}

	return grid
}

func externalLinks(r compton.Registrar, property string, links []string) compton.Element {
	linksHrefs := make(map[string]string)
	for _, link := range links {
		if linkProperty, value, ok := strings.Cut(link, "="); ok {
			linkPropertyTitle := compton_data.PropertyTitles[linkProperty]
			linksHrefs[linkPropertyTitle] = value
		}
	}
	propertyTitle := compton_data.PropertyTitles[property]
	tv := title_values.TitleValues(r, propertyTitle).AppendLinkValues(linksHrefs)
	tv.AddClass("external")
	return tv
}

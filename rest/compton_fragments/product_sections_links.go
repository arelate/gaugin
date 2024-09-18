package compton_fragments

import (
	"github.com/arelate/gaugin/rest/compton_data"
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/elements/nav_links"
)

func ProductSectionsLinks(r compton.Registrar, sections []string) compton.Element {

	productSectionsLinks := make(map[string]string)
	for _, section := range sections {
		title := compton_data.SectionTitles[section]
		productSectionsLinks[title] = "#" + title
	}

	targets := nav_links.TextLinks(
		productSectionsLinks,
		"",
		compton_data.SectionsTitlesOrder...)

	return nav_links.NavLinksTargets(r, targets...)
}

package view_models

import (
	"github.com/arelate/vangogh_local_data"
	"html/template"
)

type description struct {
	Context                string
	Description            template.HTML
	AdditionalRequirements template.HTML
	Copyrights             template.HTML
}

func NewDescription(rdx map[string][]string) *description {
	dvm := &description{Context: "iframe"}

	//Description content preparation includes the following steps:
	//1) combining DescriptionOverview and DescriptionFeatures
	//2) replacing implicit list in DescriptionFeatures with explicit list
	//3) rewriting https://items.gog.com/... links to gaugin
	//4) rewriting https://www.gog.com/game/... and https://www.gog.com/en/game/... links to gaugin
	//5) rewriting links <a href="..."/> as <a target='_top' href="..."/> to do top level navigation
	//6) fix quotes used for links in some products

	desc := propertyFromRedux(rdx, vangogh_local_data.DescriptionOverviewProperty)
	desc += implicitToExplicitList(propertyFromRedux(rdx, vangogh_local_data.DescriptionFeaturesProperty))

	desc = rewriteItemsLinks(desc)
	desc = rewriteGameLinks(desc)
	desc = rewriteLinksAsTargetTop(desc)
	desc = fixQuotes(desc)

	dvm.Description = template.HTML(desc)

	if ar := propertyFromRedux(rdx, vangogh_local_data.AdditionalRequirementsProperty); ar != "" {
		dvm.AdditionalRequirements = template.HTML(ar)
	}

	if c := propertyFromRedux(rdx, vangogh_local_data.CopyrightsProperty); c != "" {
		dvm.Copyrights = template.HTML(c)
	}

	return dvm
}

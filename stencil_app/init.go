package stencil_app

import (
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/stencil"
)

const (
	appTitle       = "gaugin"
	appAccentColor = "blueviolet"
)

func Init() (*stencil.AppConfiguration, error) {

	app := stencil.NewAppConfig(appTitle, appAccentColor)

	app.SetNavigation(NavItems, NavIcons, NavHrefs)
	app.SetFooter(FooterLocation, FooterRepoUrl)

	if err := app.SetCommonConfiguration(
		Labels,
		HiddenLabels,
		Icons,
		vangogh_local_data.TitleProperty,
		PropertyTitles,
		SectionTitles,
		DigestTitles,
		nil); err != nil {
		return app, nil
	}

	if err := app.SetListConfiguration(
		ProductsProperties,
		ProductsHiddenProperties,
		ProductPath,
		vangogh_local_data.VerticalImageProperty,
		ImagePath,
		nil); err != nil {
		return app, err
	}

	if err := app.SetItemConfiguration(
		ProductProperties,
		ProductComputedProperties,
		ProductHiddenPropertied,
		ProductSections,
		vangogh_local_data.ImageProperty,
		ImagePath,
		nil); err != nil {
		return app, err
	}

	app.SetFormatterConfiguration(fmtTitle, fmtHref, fmtClass, nil)

	if err := app.SetSearchConfiguration(
		SearchProperties,
		SearchScopes,
		SearchScopeQueries()); err != nil {
		return app, err
	}

	return app, nil

}

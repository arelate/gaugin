package stencil_app

import (
	"fmt"
	"github.com/arelate/gaugin/data"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/stencil"
	"strconv"
	"strings"
)

const (
	appTitle       = "gaugin"
	appAccentColor = "blueviolet"
)

func Init() (*stencil.App, error) {

	app := stencil.NewApp(appTitle, appAccentColor)

	app.SetNavigation(NavItems, NavIcons, NavHrefs)
	app.SetFooter(FooterLocation, FooterRepoUrl)

	app.SetTitles(vangogh_local_data.TitleProperty, PropertyTitles, SectionTitles, DigestTitles)

	if err := app.SetLabels(ProductsLabels, nil); err != nil {
		return app, err
	}
	if err := app.SetIcons(Icons, nil); err != nil {
		return app, err
	}

	app.SetLinkParams(ProductPath, ImagePath, fmtTitle, fmtHref, fmtClass)

	if err := app.SetListParams(
		vangogh_local_data.VerticalImageProperty,
		ProductsProperties,
		ProductsClassProperties,
		nil); err != nil {
		return app, err
	}

	//if err := app.SetItemParams(BookProperties, BookSections); err != nil {
	//	return app, err
	//}

	if err := app.SetSearchParams(SearchScopes, SearchScopeQueries(), SearchProperties); err != nil {
		return app, err
	}

	return app, nil

}

var labelTitles = map[string]string{
	vangogh_local_data.OwnedProperty:          "Own",
	vangogh_local_data.TBAProperty:            "TBA",
	vangogh_local_data.ComingSoonProperty:     "Soon",
	vangogh_local_data.PreOrderProperty:       "PO",
	vangogh_local_data.InDevelopmentProperty:  "In Dev",
	vangogh_local_data.IsUsingDOSBoxProperty:  "DOSBox",
	vangogh_local_data.IsUsingScummVMProperty: "ScummVM",
	vangogh_local_data.IsFreeProperty:         "Free",
	vangogh_local_data.WishlistedProperty:     "Wish",
}

const (
	transitiveOpen  = " ("
	transitiveClose = ")"
)

func transitiveDst(s string) string {
	dst := s
	if strings.Contains(s, transitiveOpen) {
		dst = s[:strings.LastIndex(s, transitiveOpen)]
	}
	return dst
}

func transitiveSrc(s string) string {
	src := ""
	if strings.Contains(s, transitiveOpen) {
		from, to := strings.LastIndex(s, transitiveOpen)+len(transitiveOpen), strings.Index(s, transitiveClose)
		if from > to {
			to = strings.LastIndex(s, transitiveClose)
			if from > to {
				from = 0
				to = len(s) - 1
			}
		}
		src = s[from:to]
	}
	return src
}

func discountPercentageLabel(value string) string {
	dp, dl := 0, ""
	if dpi, err := strconv.Atoi(value); err == nil {
		dp = dpi
		if dp >= 80 {
			dl = "\u2158" // 4/5
		} else if dp >= 75 {
			dl = "\u00be" // 3/4
		} else if dp >= 66 {
			dl = "\u2154" // 2/3
		} else if dp >= 60 {
			dl = "\u2157" // 3/5
		} else if dp >= 50 {
			dl = "\u00bd" // 1/2
		} else if dp >= 40 {
			dl = "\u2156" // 2/5
		} else if dp >= 33 {
			dl = "\u2153" // 1/3
		} else if dp >= 25 {
			dl = "\u00bc" // 1/4
		} else if dp >= 20 {
			dl = "\u2155" // 1/5
		}
	}
	return dl
}

func ownedValidationResult(id string, rxa kvas.ReduxAssets) string {
	vr, _ := rxa.GetFirstVal(vangogh_local_data.ValidationResultProperty, id)
	return vr
}

func fmtClass(id, property, link string, rxa kvas.ReduxAssets) string {
	switch property {
	case vangogh_local_data.OwnedProperty:
		return ownedValidationResult(id, rxa)
	}
	return ""
}

func fmtHref(id, property, link string, rxa kvas.ReduxAssets) string {
	switch property {
	case vangogh_local_data.GOGOrderDateProperty:
		//FIXME
		//link = justTheDate(link)
	case vangogh_local_data.PublishersProperty:
		fallthrough
	case vangogh_local_data.DevelopersProperty:
		return fmt.Sprintf("/search?%s=%s&sort=global-release-date&desc=true", property, link)
	case vangogh_local_data.IncludesGamesProperty:
		fallthrough
	case vangogh_local_data.IsIncludedByGamesProperty:
		fallthrough
	case vangogh_local_data.RequiresGamesProperty:
		fallthrough
	case vangogh_local_data.IsRequiredByGamesProperty:
		return fmt.Sprintf("/product?id=%s", transitiveSrc(link))
	case data.GauginGOGLinksProperty:
		//FIXME
		//return gogLink(transitiveSrc(link))
	case data.GauginSteamLinksProperty:
		return transitiveSrc(link)
	}
	return fmt.Sprintf("/search?%s=%s", property, link)
}

func fmtTitle(id, property, link string, rxa kvas.ReduxAssets) string {
	title := link

	switch property {
	case vangogh_local_data.WishlistedProperty:
		fallthrough
	case vangogh_local_data.OwnedProperty:
		fallthrough
	case vangogh_local_data.PreOrderProperty:
		fallthrough
	case vangogh_local_data.ComingSoonProperty:
		fallthrough
	case vangogh_local_data.TBAProperty:
		fallthrough
	case vangogh_local_data.InDevelopmentProperty:
		fallthrough
	case vangogh_local_data.IsUsingDOSBoxProperty:
		fallthrough
	case vangogh_local_data.IsUsingScummVMProperty:
		fallthrough
	case vangogh_local_data.IsFreeProperty:
		if link == "true" {
			return labelTitles[property]
		}
		return ""
	case vangogh_local_data.ProductTypeProperty:
		if link == "GAME" {
			return ""
		}
	case vangogh_local_data.DiscountPercentageProperty:
		if link == "0" {
			return ""
		}
		return "Sale " + discountPercentageLabel(link)
	case vangogh_local_data.TagIdProperty:
		return transitiveDst(link)
	case vangogh_local_data.IncludesGamesProperty:
		fallthrough
	case vangogh_local_data.IsIncludedByGamesProperty:
		fallthrough
	case vangogh_local_data.RequiresGamesProperty:
		fallthrough
	case vangogh_local_data.IsRequiredByGamesProperty:
		title = transitiveDst(link)
	case vangogh_local_data.GOGOrderDateProperty:
		//FIXME
		//title = justTheDate(link)
	case vangogh_local_data.LanguageCodeProperty:
		//FIXME
		//title = languageCodeFlag(transitiveSrc(link)) + " " + transitiveDst(link)
	case data.GauginGOGLinksProperty:
		fallthrough
	case data.GauginSteamLinksProperty:
		title = PropertyTitles[transitiveDst(link)]
	}

	return title
}

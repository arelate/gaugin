package view_models

import (
	"fmt"
	"github.com/arelate/gaugin/data"
	"github.com/arelate/gaugin/stencil_app"
	"github.com/arelate/gog_integration"
	"github.com/arelate/steam_integration"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/issa"
	"html/template"
	"net/url"
	"strconv"
	"strings"
)

type Product struct {
	Context string
	Id      string
	// Title
	Title string
	//Image
	DehydratedImage template.URL
	Image           string
	// Labels
	Labels *labels
	// Special format properties
	BasePrice        string
	Price            string
	OperatingSystems []string
	Rating           string
	// Text properties
	Properties      map[string]map[string]string
	PropertyOrder   []string
	PropertyTitles  map[string]string
	PropertyClasses map[string]string
	// Sections
	Sections      []string
	SectionTitles map[string]string
}

func NewProduct(redux vangogh_local_data.IdReduxAssets) (*Product, error) {
	switch len(redux) {
	case 0:
		return nil, fmt.Errorf("empty rdx")
	case 1:
		for id, rdx := range redux {

			pvm := &Product{
				Context: "product",
				Id:      id,

				DehydratedImage: template.URL(
					issa.Hydrate(
						propertyFromRedux(rdx, vangogh_local_data.DehydratedImageProperty))),
				Image: propertyFromRedux(rdx, vangogh_local_data.ImageProperty),
				Title: propertyFromRedux(rdx, vangogh_local_data.TitleProperty),

				Labels: NewLabels(rdx),

				BasePrice: propertyFromRedux(rdx, vangogh_local_data.BasePriceProperty),
				Price:     propertyFromRedux(rdx, vangogh_local_data.PriceProperty),

				OperatingSystems: propertiesFromRedux(rdx, vangogh_local_data.OperatingSystemsProperty),
				Rating:           propertyFromRedux(rdx, vangogh_local_data.RatingProperty),

				Properties:      make(map[string]map[string]string),
				PropertyOrder:   detailsPropertyOrder,
				PropertyTitles:  stencil_app.PropertyTitles,
				PropertyClasses: make(map[string]string),

				Sections:      make([]string, 0),
				SectionTitles: stencil_app.SectionTitles,
			}

			for _, p := range []string{
				vangogh_local_data.StoreUrlProperty,
				vangogh_local_data.ForumUrlProperty,
				vangogh_local_data.SupportUrlProperty} {
				rdx[data.GauginGOGLinksProperty] = append(rdx[data.GauginGOGLinksProperty],
					fmt.Sprintf("%s (%s)", p, propertyFromRedux(rdx, p)))
			}

			steamAppId := propertyFromRedux(rdx, vangogh_local_data.SteamAppIdProperty)
			if steamAppId != "" {
				if appId, err := strconv.ParseUint(steamAppId, 10, 32); err == nil {
					if scu := steam_integration.SteamCommunityUrl(uint32(appId)); scu != nil {
						rdx[data.GauginSteamLinksProperty] = append(rdx[data.GauginSteamLinksProperty],
							fmt.Sprintf("%s (%s)", data.GauginSteamCommunityUrlProperty, scu.String()))
					}
				}
			}

			for _, lp := range detailsPropertyOrder {
				pvm.Properties[lp] = getPropertyLinks(lp, rdx)
			}

			for _, cp := range propertyClasses {
				pvm.PropertyClasses[cp] = getPropertyClass(cp, rdx)
			}

			return pvm, nil
		}
	default:
		return nil, fmt.Errorf("too many ids, rdx")
	}
	return nil, nil
}

func getPropertyLinks(property string, rdx map[string][]string) map[string]string {

	propertyLinks := make(map[string]string)

	for _, value := range propertiesFromRedux(rdx, property) {

		linkTitle := formatPropertyLinkTitle(property, value)
		propertyLinks[linkTitle] = formatPropertyLinkHref(property, value)
	}

	return propertyLinks
}

func formatPropertyLinkTitle(property, link string) string {
	title := link

	switch property {
	case vangogh_local_data.IncludesGamesProperty:
		fallthrough
	case vangogh_local_data.IsIncludedByGamesProperty:
		fallthrough
	case vangogh_local_data.RequiresGamesProperty:
		fallthrough
	case vangogh_local_data.IsRequiredByGamesProperty:
		title = transitiveDst(link)
	case vangogh_local_data.GOGOrderDateProperty:
		title = justTheDate(link)
	case vangogh_local_data.LanguageCodeProperty:
		title = languageCodeFlag(transitiveSrc(link)) + " " + transitiveDst(link)
	case data.GauginGOGLinksProperty:
		fallthrough
	case data.GauginSteamLinksProperty:
		title = stencil_app.PropertyTitles[transitiveDst(link)]
	}

	return title
}

func formatPropertyLinkHref(property, link string) string {
	switch property {
	case vangogh_local_data.GOGOrderDateProperty:
		link = justTheDate(link)
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
		return gogLink(transitiveSrc(link))
	case data.GauginSteamLinksProperty:
		return transitiveSrc(link)
	}
	return fmt.Sprintf("/search?%s=%s", property, link)
}

func justTheDate(s string) string {
	return strings.Split(s, " ")[0]
}

func gogLink(p string) string {
	u := url.URL{
		Scheme: gog_integration.HttpsScheme,
		Host:   gog_integration.WwwGogHost,
		Path:   p,
	}
	return u.String()
}

func steamReviewClass(sr string) string {
	if strings.Contains(sr, "Positive") {
		return "positive"
	} else if strings.Contains(sr, "Negative") {
		return "negative"
	} else {
		return "neutral"
	}
}

func getPropertyClass(property string, rdx map[string][]string) string {
	switch property {
	case vangogh_local_data.SteamReviewScoreDescProperty:
		return steamReviewClass(propertyFromRedux(rdx, vangogh_local_data.SteamReviewScoreDescProperty))
	}

	return ""
}

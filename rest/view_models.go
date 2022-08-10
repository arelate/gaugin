package rest

import (
	"fmt"
	"github.com/arelate/steam_integration"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/issa"
	"html/template"
	"sort"
	"strconv"
	"strings"
)

type labels struct {
	Wishlisted         bool
	Owned              bool
	PreOrder           bool
	InDevelopment      bool
	TBA                bool
	ComingSoon         bool
	IsUsingDOSBox      bool
	IsUsingScummVM     bool
	Free               bool
	Discounted         bool
	DiscountPercentage int
	DiscountLabel      string
	Tags               []string
	LocalTags          []string
	ProductType        string
}

type listProductViewModel struct {
	Id               string
	Title            string
	Developers       []string
	Publishers       []string
	Labels           *labels
	OperatingSystems []string
}

type listViewModel struct {
	Context  string
	Products []listProductViewModel
}

type updatesViewModel struct {
	Context         string
	Sections        []string
	SectionProducts map[string]*listViewModel
}

type productDownloads struct {
	CurrentOS        bool
	OperatingSystems string
	Installers       vangogh_local_data.DownloadsList
	DLCs             vangogh_local_data.DownloadsList
	Extras           vangogh_local_data.DownloadsList
}

var gauginPropertyOrder = []string{
	vangogh_local_data.SteamReviewScoreDescProperty,
	vangogh_local_data.DevelopersProperty,
	vangogh_local_data.PublishersProperty,
	vangogh_local_data.SeriesProperty,
	vangogh_local_data.GenresProperty,
	vangogh_local_data.StoreTagsProperty,
	vangogh_local_data.SteamTagsProperty,
	vangogh_local_data.FeaturesProperty,
	vangogh_local_data.LanguageCodeProperty,
	vangogh_local_data.GlobalReleaseDateProperty,
	vangogh_local_data.GOGReleaseDateProperty,
	vangogh_local_data.GOGOrderDateProperty,
	vangogh_local_data.IncludesGamesProperty,
	vangogh_local_data.IsIncludedByGamesProperty,
	vangogh_local_data.RequiresGamesProperty,
	vangogh_local_data.IsRequiredByGamesProperty,
	GauginGOGLinksProperty,
	GauginSteamLinksProperty,
}

const (
	GauginGOGLinksProperty          = "gog-links"
	GauginSteamLinksProperty        = "steam-links"
	GauginSteamCommunityUrlProperty = "steam-community-url"
)

type productViewModel struct {
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
	Properties     map[string]map[string]string
	PropertyOrder  []string
	PropertyTitles map[string]string
	// Sections
	Sections      []string
	SectionTitles map[string]string
}

type descriptionViewModel struct {
	Context                string
	Description            template.HTML
	AdditionalRequirements template.HTML
	Copyrights             template.HTML
}

type changelogViewModel struct {
	Context   string
	Changelog template.HTML
}

type screenshotsViewModel struct {
	Context     string
	Screenshots []string
}

type videosViewModel struct {
	Context      string
	LocalVideos  []string
	RemoteVideos []string
}

type downloadsViewModel struct {
	Context   string
	CurrentOS *productDownloads
	OtherOS   *productDownloads
	Extras    *productDownloads
}

type newsItemViewModel struct {
	Title     string
	Date      int
	Author    string
	Url       string
	Tags      string
	FeedLabel string
	Contents  template.HTML
}

func steamAppNewsViewModelFromResponse(san *steam_integration.AppNews) *steamAppNewsViewModel {
	sanvm := &steamAppNewsViewModel{
		Context:   "iframe",
		Count:     san.Count,
		NewsItems: make([]*newsItemViewModel, 0, len(san.NewsItems)),
	}

	for _, ni := range san.NewsItems {
		sanvm.NewsItems = append(sanvm.NewsItems, &newsItemViewModel{
			Title:     ni.Title,
			Date:      ni.Date,
			Author:    ni.Author,
			Url:       ni.Url,
			Tags:      strings.Join(ni.Tags, ","),
			FeedLabel: ni.FeedLabel,
			Contents:  template.HTML(steamAppNewsToHTML(ni.Contents)),
		})
	}

	return sanvm
}

type steamAppNewsViewModel struct {
	Context   string
	Count     uint32
	NewsItems []*newsItemViewModel
}

type tagsEditViewModel struct {
	Context      string
	Id           string
	Title        string
	Owned        bool
	AllTags      []string
	AllLocalTags []string
	Tags         map[string]bool
	LocalTags    map[string]bool
}

func propertyFromRedux(redux map[string][]string, property string) string {
	properties := propertiesFromRedux(redux, property)
	if len(properties) > 0 {
		return properties[0]
	}
	return ""
}

func flagFromRedux(redux map[string][]string, property string) bool {
	return propertyFromRedux(redux, property) == vangogh_local_data.TrueValue
}

func propertiesFromRedux(redux map[string][]string, property string) []string {
	if val, ok := redux[property]; ok {
		return val
	} else {
		return []string{}
	}
}

func discountPercentageLabelFromRedux(redux map[string][]string) (int, string) {
	dp, dl := 0, ""
	dpa := propertyFromRedux(redux, vangogh_local_data.DiscountPercentageProperty)
	if dpi, err := strconv.Atoi(dpa); err == nil {
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
	return dp, dl
}

func listViewModelFromRedux(order []string, redux map[string]map[string][]string) *listViewModel {
	lvm := &listViewModel{
		Products: make([]listProductViewModel, 0, len(order)),
	}
	for _, id := range order {
		rdx, ok := redux[id]
		if !ok {
			continue
		}
		lpvm := listProductViewModel{
			Id:               id,
			Title:            propertyFromRedux(rdx, vangogh_local_data.TitleProperty),
			Labels:           labelsFromRedux(rdx),
			Developers:       propertiesFromRedux(rdx, vangogh_local_data.DevelopersProperty),
			Publishers:       propertiesFromRedux(rdx, vangogh_local_data.PublishersProperty),
			OperatingSystems: propertiesFromRedux(rdx, vangogh_local_data.OperatingSystemsProperty),
		}

		lvm.Products = append(lvm.Products, lpvm)
	}
	return lvm
}

func labelsFromRedux(rdx map[string][]string) *labels {
	lbs := &labels{
		Wishlisted:     flagFromRedux(rdx, vangogh_local_data.WishlistedProperty),
		Owned:          flagFromRedux(rdx, vangogh_local_data.OwnedProperty),
		Free:           flagFromRedux(rdx, vangogh_local_data.IsFreeProperty),
		Discounted:     flagFromRedux(rdx, vangogh_local_data.IsDiscountedProperty),
		PreOrder:       flagFromRedux(rdx, vangogh_local_data.PreOrderProperty),
		TBA:            flagFromRedux(rdx, vangogh_local_data.TBAProperty),
		ComingSoon:     flagFromRedux(rdx, vangogh_local_data.ComingSoonProperty),
		InDevelopment:  flagFromRedux(rdx, vangogh_local_data.InDevelopmentProperty),
		IsUsingDOSBox:  flagFromRedux(rdx, vangogh_local_data.IsUsingDOSBoxProperty),
		IsUsingScummVM: flagFromRedux(rdx, vangogh_local_data.IsUsingScummVMProperty),
		Tags:           propertiesFromRedux(rdx, vangogh_local_data.TagIdProperty),
		LocalTags:      propertiesFromRedux(rdx, vangogh_local_data.LocalTagsProperty),
		ProductType:    propertyFromRedux(rdx, vangogh_local_data.ProductTypeProperty),
	}

	lbs.DiscountPercentage, lbs.DiscountLabel = discountPercentageLabelFromRedux(rdx)

	return lbs
}

func productViewModelFromRedux(redux map[string]map[string][]string) (*productViewModel, error) {
	switch len(redux) {
	case 0:
		return nil, fmt.Errorf("empty rdx")
	case 1:
		for id, rdx := range redux {

			pvm := &productViewModel{
				Context: "product",
				Id:      id,

				DehydratedImage: template.URL(
					issa.Hydrate(
						propertyFromRedux(rdx, vangogh_local_data.DehydratedImageProperty))),
				Image: propertyFromRedux(rdx, vangogh_local_data.ImageProperty),
				Title: propertyFromRedux(rdx, vangogh_local_data.TitleProperty),

				Labels: labelsFromRedux(rdx),

				BasePrice: propertyFromRedux(rdx, vangogh_local_data.BasePriceProperty),
				Price:     propertyFromRedux(rdx, vangogh_local_data.PriceProperty),

				OperatingSystems: propertiesFromRedux(rdx, vangogh_local_data.OperatingSystemsProperty),
				Rating:           propertyFromRedux(rdx, vangogh_local_data.RatingProperty),

				Properties:     make(map[string]map[string]string),
				PropertyOrder:  gauginPropertyOrder,
				PropertyTitles: propertyTitles,

				Sections:      make([]string, 0),
				SectionTitles: sectionTitles,
			}

			for _, p := range []string{
				vangogh_local_data.StoreUrlProperty,
				vangogh_local_data.ForumUrlProperty,
				vangogh_local_data.SupportUrlProperty} {
				rdx[GauginGOGLinksProperty] = append(rdx[GauginGOGLinksProperty],
					fmt.Sprintf("%s (%s)", p, propertyFromRedux(rdx, p)))
			}

			steamAppId := propertyFromRedux(rdx, vangogh_local_data.SteamAppIdProperty)
			if steamAppId != "" {
				if appId, err := strconv.Atoi(steamAppId); err == nil {
					if scu := steam_integration.SteamCommunityUrl(uint32(appId)); scu != nil {
						rdx[GauginSteamLinksProperty] = append(rdx[GauginSteamLinksProperty],
							fmt.Sprintf("%s (%s)", GauginSteamCommunityUrlProperty, scu.String()))
					}
				}
			}

			for _, lp := range gauginPropertyOrder {
				pvm.Properties[lp] = getPropertyLinks(lp, rdx)
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
	case vangogh_local_data.GlobalReleaseDateProperty:
		fallthrough
	case vangogh_local_data.GOGReleaseDateProperty:
		title = formatDate(link)
	case vangogh_local_data.GOGOrderDateProperty:
		title = formatDate(justTheDate(link))
	case vangogh_local_data.LanguageCodeProperty:
		title = languageCodeFlag(transitiveSrc(link)) + " " + transitiveDst(link)
	case GauginGOGLinksProperty:
		fallthrough
	case GauginSteamLinksProperty:
		title = propertyTitles[transitiveDst(link)]
	}

	return title
}

func formatPropertyLinkHref(property, link string) string {
	switch property {
	case vangogh_local_data.GlobalReleaseDateProperty:
		fallthrough
	case vangogh_local_data.GOGReleaseDateProperty:
		fallthrough
	case vangogh_local_data.GOGOrderDateProperty:
		return ""
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
	case GauginGOGLinksProperty:
		return gogLink(transitiveSrc(link))
	case GauginSteamLinksProperty:
		return transitiveSrc(link)
	default:
		return fmt.Sprintf("/search?%s=%s", property, link)
	}

	return ""
}

func updatesViewModelFromRedux(
	updates map[string][]string,
	rdx map[string]map[string][]string) *updatesViewModel {

	sections := make([]string, 0, len(updates))
	sectionProducts := make(map[string]*listViewModel)
	for s, ids := range updates {
		sections = append(sections, s)
		sectionProducts[s] = listViewModelFromRedux(ids, rdx)
	}

	sort.Strings(sections)

	uvm := &updatesViewModel{
		Context:         "updates",
		Sections:        sections,
		SectionProducts: sectionProducts,
	}

	return uvm
}

package rest

import (
	"fmt"
	"github.com/arelate/steam_integration"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/issa"
	"html/template"
	"sort"
	"strconv"
)

type listProductViewModel struct {
	Id                 string
	Title              string
	Developers         []string
	Publisher          string
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
	OperatingSystems   []string
	Tags               []string
	LocalTags          []string
	ProductType        string
}

type listViewModel struct {
	Context  string
	Products []listProductViewModel
}

type updatesViewModel struct {
	Since           int
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

type productViewModel struct {
	Context string
	Id      string
	// text properties
	ProductType       string
	Title             string
	DehydratedImage   template.URL
	Image             string
	Tags              []string
	LocalTags         []string
	SteamTags         []string
	OperatingSystems  []string
	Rating            string
	Developers        []string
	Publisher         string
	Series            string
	Genres            []string
	Properties        []string
	Features          []string
	LanguageCodes     []string
	GlobalReleaseDate string
	GOGReleaseDate    string
	GOGOrderDate      string
	IncludesGames     []string
	IsIncludedByGames []string
	RequiresGames     []string
	IsRequiredByGames []string
	// urls
	StoreUrl   string
	ForumUrl   string
	SupportUrl string
	// video-ids
	Videos []string
	// downloads
	CurrentOS     *productDownloads
	OtherOS       *productDownloads
	LanguageFlags map[string]string
	// labels
	Wishlisted     bool
	Owned          bool
	Free           bool
	Discounted     bool
	PreOrder       bool
	TBA            bool
	ComingSoon     bool
	InDevelopment  bool
	IsUsingDOSBox  bool
	IsUsingScummVM bool
	// price
	BasePrice          string
	Price              string
	DiscountPercentage int
	DiscountLabel      string
	// Steam Community url
	SteamCommunityUrl    string
	SteamAppId           string
	SteamReviewScoreDesc string
	// has properties and data
	HasDescription  bool
	HasChangelog    bool
	HasScreenshots  bool
	HasSteamAppNews bool
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

type newsItemViewModel struct {
	Title    string
	Date     int
	Author   string
	Url      string
	Contents template.HTML
}

func steamAppNewsViewModelFromResponse(san *steam_integration.AppNews) *steamAppNewsViewModel {
	sanvm := &steamAppNewsViewModel{
		Context:   "iframe",
		Count:     san.Count,
		NewsItems: make([]*newsItemViewModel, 0, len(san.NewsItems)),
	}

	for _, ni := range san.NewsItems {
		sanvm.NewsItems = append(sanvm.NewsItems, &newsItemViewModel{
			Title:    ni.Title,
			Date:     ni.Date,
			Author:   ni.Author,
			Url:      ni.Url,
			Contents: template.HTML(steamAppNewsToHTML(ni.Contents)),
		})
	}

	return sanvm
}

type steamAppNewsViewModel struct {
	Context   string
	Count     uint32
	NewsItems []*newsItemViewModel
}

func propertyFromRedux(redux map[string][]string, property string) string {
	properties := propertiesFromRedux(redux, property)
	if len(properties) > 0 {
		return properties[0]
	}
	return ""
}

func flagFromRedux(redux map[string][]string, property string) bool {
	return propertyFromRedux(redux, property) == "true"
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
			Wishlisted:       flagFromRedux(rdx, vangogh_local_data.WishlistedProperty),
			Owned:            flagFromRedux(rdx, vangogh_local_data.OwnedProperty),
			Free:             flagFromRedux(rdx, vangogh_local_data.IsFreeProperty),
			Discounted:       flagFromRedux(rdx, vangogh_local_data.IsDiscountedProperty),
			PreOrder:         flagFromRedux(rdx, vangogh_local_data.PreOrderProperty),
			ComingSoon:       flagFromRedux(rdx, vangogh_local_data.ComingSoonProperty),
			InDevelopment:    flagFromRedux(rdx, vangogh_local_data.InDevelopmentProperty),
			TBA:              flagFromRedux(rdx, vangogh_local_data.TBAProperty),
			IsUsingDOSBox:    flagFromRedux(rdx, vangogh_local_data.IsUsingDOSBoxProperty),
			IsUsingScummVM:   flagFromRedux(rdx, vangogh_local_data.IsUsingScummVMProperty),
			Developers:       propertiesFromRedux(rdx, vangogh_local_data.DevelopersProperty),
			Publisher:        propertyFromRedux(rdx, vangogh_local_data.PublisherProperty),
			OperatingSystems: propertiesFromRedux(rdx, vangogh_local_data.OperatingSystemsProperty),
			Tags:             propertiesFromRedux(rdx, vangogh_local_data.TagIdProperty),
			LocalTags:        propertiesFromRedux(rdx, vangogh_local_data.LocalTagsProperty),
			ProductType:      propertyFromRedux(rdx, vangogh_local_data.ProductTypeProperty),
		}

		lpvm.DiscountPercentage, lpvm.DiscountLabel = discountPercentageLabelFromRedux(rdx)

		lvm.Products = append(lvm.Products, lpvm)
	}
	return lvm
}

func productViewModelFromRedux(redux map[string]map[string][]string) (*productViewModel, error) {
	switch len(redux) {
	case 0:
		return nil, fmt.Errorf("empty rdx")
	case 1:
		for id, rdx := range redux {

			pvm := &productViewModel{
				Context:              "product",
				Id:                   id,
				DehydratedImage:      template.URL(issa.Hydrate(propertyFromRedux(rdx, vangogh_local_data.DehydratedImageProperty))),
				Image:                propertyFromRedux(rdx, vangogh_local_data.ImageProperty),
				ProductType:          propertyFromRedux(rdx, vangogh_local_data.ProductTypeProperty),
				Title:                propertyFromRedux(rdx, vangogh_local_data.TitleProperty),
				Tags:                 propertiesFromRedux(rdx, vangogh_local_data.TagIdProperty),
				LocalTags:            propertiesFromRedux(rdx, vangogh_local_data.LocalTagsProperty),
				SteamTags:            propertiesFromRedux(rdx, vangogh_local_data.SteamTagsProperty),
				OperatingSystems:     propertiesFromRedux(rdx, vangogh_local_data.OperatingSystemsProperty),
				Rating:               propertyFromRedux(rdx, vangogh_local_data.RatingProperty),
				Developers:           propertiesFromRedux(rdx, vangogh_local_data.DevelopersProperty),
				Publisher:            propertyFromRedux(rdx, vangogh_local_data.PublisherProperty),
				Series:               propertyFromRedux(rdx, vangogh_local_data.SeriesProperty),
				Genres:               propertiesFromRedux(rdx, vangogh_local_data.GenresProperty),
				Properties:           propertiesFromRedux(rdx, vangogh_local_data.PropertiesProperty),
				Features:             propertiesFromRedux(rdx, vangogh_local_data.FeaturesProperty),
				LanguageCodes:        propertiesFromRedux(rdx, vangogh_local_data.LanguageCodeProperty),
				GlobalReleaseDate:    propertyFromRedux(rdx, vangogh_local_data.GlobalReleaseDateProperty),
				GOGReleaseDate:       propertyFromRedux(rdx, vangogh_local_data.GOGReleaseDateProperty),
				GOGOrderDate:         propertyFromRedux(rdx, vangogh_local_data.GOGOrderDateProperty),
				IncludesGames:        propertiesFromRedux(rdx, vangogh_local_data.IncludesGamesProperty),
				IsIncludedByGames:    propertiesFromRedux(rdx, vangogh_local_data.IsIncludedByGamesProperty),
				RequiresGames:        propertiesFromRedux(rdx, vangogh_local_data.RequiresGamesProperty),
				IsRequiredByGames:    propertiesFromRedux(rdx, vangogh_local_data.IsRequiredByGamesProperty),
				StoreUrl:             propertyFromRedux(rdx, vangogh_local_data.StoreUrlProperty),
				ForumUrl:             propertyFromRedux(rdx, vangogh_local_data.ForumUrlProperty),
				SupportUrl:           propertyFromRedux(rdx, vangogh_local_data.SupportUrlProperty),
				Videos:               propertiesFromRedux(rdx, vangogh_local_data.VideoIdProperty),
				Wishlisted:           flagFromRedux(rdx, vangogh_local_data.WishlistedProperty),
				Owned:                flagFromRedux(rdx, vangogh_local_data.OwnedProperty),
				Free:                 flagFromRedux(rdx, vangogh_local_data.IsFreeProperty),
				Discounted:           flagFromRedux(rdx, vangogh_local_data.IsDiscountedProperty),
				PreOrder:             flagFromRedux(rdx, vangogh_local_data.PreOrderProperty),
				TBA:                  flagFromRedux(rdx, vangogh_local_data.TBAProperty),
				ComingSoon:           flagFromRedux(rdx, vangogh_local_data.ComingSoonProperty),
				InDevelopment:        flagFromRedux(rdx, vangogh_local_data.InDevelopmentProperty),
				IsUsingDOSBox:        flagFromRedux(rdx, vangogh_local_data.IsUsingDOSBoxProperty),
				IsUsingScummVM:       flagFromRedux(rdx, vangogh_local_data.IsUsingScummVMProperty),
				BasePrice:            propertyFromRedux(rdx, vangogh_local_data.BasePriceProperty),
				Price:                propertyFromRedux(rdx, vangogh_local_data.PriceProperty),
				SteamAppId:           propertyFromRedux(rdx, vangogh_local_data.SteamAppIdProperty),
				SteamReviewScoreDesc: propertyFromRedux(rdx, vangogh_local_data.SteamReviewScoreDescProperty),
			}

			pvm.DiscountPercentage, pvm.DiscountLabel = discountPercentageLabelFromRedux(rdx)

			if pvm.SteamAppId != "" {
				if appId, err := strconv.Atoi(pvm.SteamAppId); err == nil {
					if scu := steam_integration.SteamCommunityUrl(uint32(appId)); scu != nil {
						pvm.SteamCommunityUrl = scu.String()
					}
				}
			}

			return pvm, nil
		}
	default:
		return nil, fmt.Errorf("too many ids, rdx")
	}
	return nil, nil
}

func updatesViewModelFromRedux(updates map[string][]string,
	since int,
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
		Since:           since,
		Sections:        sections,
		SectionProducts: sectionProducts,
	}

	return uvm
}

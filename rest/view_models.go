package rest

import (
	"fmt"
	"github.com/arelate/vangogh_local_data"
)

type listProductViewModel struct {
	Id         string
	Title      string
	Developers []string
	Publisher  string
}

type listViewModel struct {
	Context  string
	Products []listProductViewModel
}

type productViewModel struct {
	Id         string
	Wishlisted bool
	// text properties
	Title             string
	Image             string
	Tags              []string
	OperatingSystems  []string
	Rating            string
	Developers        []string
	Publisher         string
	Series            string
	Genres            []string
	Features          []string
	LanguageCodes     []string
	GlobalReleaseDate string
	GOGReleaseDate    string
	GOGOrderDate      string
	IncludesGames     []string
	IsIncludedByGames []string
	RequiresGames     []string
	IsRequiredByGames []string
	Types             []string
	// urls
	StoreUrl   string
	ForumUrl   string
	SupportUrl string
	// long text
	Changelog   string
	Description string
	// screenshots
	Screenshots []string
	// video-ids
	Videos []string
	// downloads
	Downloads vangogh_local_data.DownloadsList
}

func propertyFromRedux(redux map[string][]string, property string) string {
	properties := propertiesFromRedux(redux, property)
	if len(properties) > 0 {
		return properties[0]
	}
	return ""
}

func propertiesFromRedux(redux map[string][]string, property string) []string {
	if val, ok := redux[property]; ok {
		return val
	} else {
		return []string{}
	}
}

func listViewModelFromRedux(order []string, redux map[string]map[string][]string) *listViewModel {
	lvm := &listViewModel{
		Products: make([]listProductViewModel, 0, len(redux)),
	}
	for _, id := range order {
		properties, ok := redux[id]
		if !ok {
			continue
		}
		lvm.Products = append(
			lvm.Products,
			listProductViewModel{
				Id:         id,
				Title:      propertyFromRedux(properties, "title"),
				Developers: propertiesFromRedux(properties, "developers"),
				Publisher:  propertyFromRedux(properties, "publisher"),
			})
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
				Id:                id,
				Image:             propertyFromRedux(rdx, vangogh_local_data.ImageProperty),
				Title:             propertyFromRedux(rdx, vangogh_local_data.TitleProperty),
				Tags:              propertiesFromRedux(rdx, vangogh_local_data.TagIdProperty),
				OperatingSystems:  propertiesFromRedux(rdx, vangogh_local_data.OperatingSystemsProperty),
				Rating:            propertyFromRedux(rdx, vangogh_local_data.RatingProperty),
				Developers:        propertiesFromRedux(rdx, vangogh_local_data.DevelopersProperty),
				Publisher:         propertyFromRedux(rdx, vangogh_local_data.PublisherProperty),
				Series:            propertyFromRedux(rdx, vangogh_local_data.SeriesProperty),
				Genres:            propertiesFromRedux(rdx, vangogh_local_data.GenresProperty),
				Features:          propertiesFromRedux(rdx, vangogh_local_data.FeaturesProperty),
				LanguageCodes:     propertiesFromRedux(rdx, vangogh_local_data.LanguageCodeProperty),
				GlobalReleaseDate: propertyFromRedux(rdx, vangogh_local_data.GlobalReleaseDateProperty),
				GOGReleaseDate:    propertyFromRedux(rdx, vangogh_local_data.GOGReleaseDateProperty),
				GOGOrderDate:      propertyFromRedux(rdx, vangogh_local_data.GOGOrderDateProperty),
				IncludesGames:     propertiesFromRedux(rdx, vangogh_local_data.IncludesGamesProperty),
				IsIncludedByGames: propertiesFromRedux(rdx, vangogh_local_data.IsIncludedByGamesProperty),
				RequiresGames:     propertiesFromRedux(rdx, vangogh_local_data.RequiresGamesProperty),
				IsRequiredByGames: propertiesFromRedux(rdx, vangogh_local_data.IsRequiredByGamesProperty),
				Types:             propertiesFromRedux(rdx, vangogh_local_data.TypesProperty),
				StoreUrl:          propertyFromRedux(rdx, vangogh_local_data.StoreUrlProperty),
				ForumUrl:          propertyFromRedux(rdx, vangogh_local_data.ForumUrlProperty),
				SupportUrl:        propertyFromRedux(rdx, vangogh_local_data.SupportUrlProperty),
				Changelog:         propertyFromRedux(rdx, vangogh_local_data.ChanglogProperty),
				Description:       propertyFromRedux(rdx, vangogh_local_data.DescriptionProperty),
				Screenshots:       propertiesFromRedux(rdx, vangogh_local_data.ScreenshotsProperty),
				Videos:            propertiesFromRedux(rdx, vangogh_local_data.VideoIdProperty),
			}

			for _, t := range pvm.Types {
				if t == vangogh_local_data.WishlistProducts.String() {
					pvm.Wishlisted = true
					break
				}
			}

			return pvm, nil
		}
	default:
		return nil, fmt.Errorf("too many ids, rdx")
	}
	return nil, nil
}

package stencil_app

import (
	"github.com/arelate/gaugin/data"
	"github.com/arelate/vangogh_local_data"
)

var ProductProperties = []string{
	vangogh_local_data.TitleProperty,

	vangogh_local_data.ImageProperty,

	vangogh_local_data.ProductTypeProperty,
	vangogh_local_data.TagIdProperty,
	vangogh_local_data.LocalTagsProperty,

	vangogh_local_data.WishlistedProperty,
	vangogh_local_data.DiscountPercentageProperty,
	vangogh_local_data.PriceProperty,
	vangogh_local_data.OperatingSystemsProperty,
	vangogh_local_data.HLTBPlatformsProperty,
	vangogh_local_data.RatingProperty,
	vangogh_local_data.SteamReviewScoreDescProperty,
	vangogh_local_data.HLTBReviewScoreProperty,
	vangogh_local_data.DevelopersProperty,
	vangogh_local_data.PublishersProperty,
	vangogh_local_data.SeriesProperty,
	vangogh_local_data.GenresProperty,
	vangogh_local_data.HLTBGenresProperty,
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

	vangogh_local_data.StoreUrlProperty,
	vangogh_local_data.ForumUrlProperty,
	vangogh_local_data.SupportUrlProperty,
	vangogh_local_data.SteamAppIdProperty,
	vangogh_local_data.PCGWPageIdProperty,
	vangogh_local_data.HLTBIdProperty,
	vangogh_local_data.IGDBIdProperty,
	vangogh_local_data.StrategyWikiIdProperty,
	vangogh_local_data.MobyGamesIdProperty,
	vangogh_local_data.WikipediaIdProperty,
	vangogh_local_data.WineHQIdProperty,
	vangogh_local_data.VNDBIdProperty,
	vangogh_local_data.IGNWikiSlugProperty,

	vangogh_local_data.HLTBHoursToCompleteMainProperty,
	vangogh_local_data.HLTBHoursToCompletePlusProperty,
	vangogh_local_data.HLTBHoursToComplete100Property,

	vangogh_local_data.OwnedProperty,
	vangogh_local_data.IsFreeProperty,
	vangogh_local_data.PreOrderProperty,
	vangogh_local_data.ComingSoonProperty,
	vangogh_local_data.InDevelopmentProperty,
	vangogh_local_data.IsUsingDOSBoxProperty,
	vangogh_local_data.IsUsingScummVMProperty,
	vangogh_local_data.ValidationResultProperty,
}

var ProductComputedProperties = []string{
	data.GauginGOGLinksProperty,
	data.GauginOtherLinksProperty,
	data.GauginSteamLinksProperty,
}

var ProductHiddenPropertied = []string{
	vangogh_local_data.OwnedProperty,
	vangogh_local_data.ValidationResultProperty,
	vangogh_local_data.ProductTypeProperty,
	vangogh_local_data.ComingSoonProperty,
	vangogh_local_data.PreOrderProperty,
	vangogh_local_data.InDevelopmentProperty,
	vangogh_local_data.IsUsingDOSBoxProperty,
	vangogh_local_data.IsUsingScummVMProperty,
	vangogh_local_data.IsFreeProperty,
	vangogh_local_data.DiscountPercentageProperty,
	vangogh_local_data.ValidationResultProperty,
	vangogh_local_data.StoreUrlProperty,
	vangogh_local_data.ForumUrlProperty,
	vangogh_local_data.SupportUrlProperty,
	vangogh_local_data.SteamAppIdProperty,
	vangogh_local_data.PCGWPageIdProperty,
	vangogh_local_data.HLTBIdProperty,
	vangogh_local_data.IGDBIdProperty,
	vangogh_local_data.StrategyWikiIdProperty,
	vangogh_local_data.MobyGamesIdProperty,
	vangogh_local_data.WikipediaIdProperty,
	vangogh_local_data.WineHQIdProperty,
	vangogh_local_data.VNDBIdProperty,
	vangogh_local_data.IGNWikiSlugProperty,
}

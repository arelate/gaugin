package compton_data

import (
	"github.com/arelate/gaugin/data"
	"github.com/arelate/vangogh_local_data"
)

var PropertyTitles = map[string]string{
	vangogh_local_data.TitleProperty:                             "Title",
	vangogh_local_data.DescriptionOverviewProperty:               "Description",
	vangogh_local_data.TagIdProperty:                             "Account Tags",
	vangogh_local_data.LocalTagsProperty:                         "Local Tags",
	vangogh_local_data.SteamTagsProperty:                         "Steam Tags",
	vangogh_local_data.OperatingSystemsProperty:                  "OS",
	vangogh_local_data.DevelopersProperty:                        "Developers",
	vangogh_local_data.PublishersProperty:                        "Publishers",
	vangogh_local_data.EnginesProperty:                           "Engine",
	vangogh_local_data.EnginesBuildsProperty:                     "Engine Build",
	vangogh_local_data.SeriesProperty:                            "Series",
	vangogh_local_data.GenresProperty:                            "Genres",
	vangogh_local_data.StoreTagsProperty:                         "Store Tags",
	vangogh_local_data.FeaturesProperty:                          "Features",
	vangogh_local_data.LanguageCodeProperty:                      "Language",
	vangogh_local_data.IncludesGamesProperty:                     "Includes",
	vangogh_local_data.IsIncludedByGamesProperty:                 "Included By",
	vangogh_local_data.RequiresGamesProperty:                     "Requires",
	vangogh_local_data.IsRequiredByGamesProperty:                 "Required By",
	vangogh_local_data.ProductTypeProperty:                       "Product Type",
	vangogh_local_data.WishlistedProperty:                        "Wishlisted",
	vangogh_local_data.OwnedProperty:                             "Owned",
	vangogh_local_data.IsFreeProperty:                            "Free",
	vangogh_local_data.IsDiscountedProperty:                      "On Sale",
	vangogh_local_data.PreOrderProperty:                          "Pre-order",
	vangogh_local_data.ComingSoonProperty:                        "Coming Soon",
	vangogh_local_data.InDevelopmentProperty:                     "In Development",
	vangogh_local_data.TypesProperty:                             "Data Type",
	vangogh_local_data.SteamReviewScoreDescProperty:              "Steam Reviews",
	vangogh_local_data.SteamDeckAppCompatibilityCategoryProperty: "Steam Deck",
	vangogh_local_data.ProtonDBTierProperty:                      "ProtonDB Tier",
	vangogh_local_data.ProtonDBConfidenceProperty:                "ProtonDB Confidence",
	vangogh_local_data.SortProperty:                              "Sort",
	vangogh_local_data.DescendingProperty:                        "Descending",
	vangogh_local_data.GlobalReleaseDateProperty:                 "Global Release",
	vangogh_local_data.GOGReleaseDateProperty:                    "GOG.com Release",
	vangogh_local_data.GOGOrderDateProperty:                      "GOG.com Order",
	vangogh_local_data.ValidationResultProperty:                  "Validation Result",
	vangogh_local_data.RatingProperty:                            "Rating",
	vangogh_local_data.PriceProperty:                             "Price",
	vangogh_local_data.BasePriceProperty:                         "Base Price",
	vangogh_local_data.DiscountPercentageProperty:                "Discount",

	vangogh_local_data.HLTBHoursToCompleteMainProperty: "HLTB Main Story",
	vangogh_local_data.HLTBHoursToCompletePlusProperty: "HLTB Story + Extras",
	vangogh_local_data.HLTBHoursToComplete100Property:  "HLTB Completionist",
	vangogh_local_data.HLTBGenresProperty:              "HLTB Genres",
	vangogh_local_data.HLTBPlatformsProperty:           "HLTB Platforms",
	vangogh_local_data.HLTBReviewScoreProperty:         "HLTB Review Score",

	data.GauginGOGLinksProperty:   "GOG.com Links",
	data.GauginOtherLinksProperty: "Other Links",
	data.GauginSteamLinksProperty: "Steam Links",

	vangogh_local_data.ForumUrlProperty:   "Forum",
	vangogh_local_data.StoreUrlProperty:   "Store",
	vangogh_local_data.SupportUrlProperty: "Support",

	data.GauginSteamCommunityUrlProperty: "Community",

	data.GauginGOGDBUrlProperty:        "GOGDB",
	data.GauginIGDBUrlProperty:         "IGDB",
	data.GauginHLTBUrlProperty:         "HLTB",
	data.GauginMobyGamesUrlProperty:    "MobyGames",
	data.GauginPCGamingWikiUrlProperty: "PCGamingWiki",
	data.GauginProtonDBUrlProperty:     "ProtonDB",
	data.GauginStrategyWikiUrlProperty: "StrategyWiki",
	data.GauginWikipediaUrlProperty:    "Wikipedia",
	data.GauginWineHQUrlProperty:       "WineHQ",
	data.GauginVNDBUrlProperty:         "VNDB",
	data.GauginIGNWikiUrlProperty:      "IGN Wiki",

	vangogh_local_data.TrueValue:  "Yes",
	vangogh_local_data.FalseValue: "No",

	vangogh_local_data.AccountProducts.String():      "Account Products",
	vangogh_local_data.ApiProductsV1.String():        "API Products v1",
	vangogh_local_data.ApiProductsV2.String():        "API Products v2",
	vangogh_local_data.CatalogProducts.String():      "Catalog Products",
	vangogh_local_data.Details.String():              "Details",
	vangogh_local_data.HLTBData.String():             "HowLongToBeat Data",
	vangogh_local_data.HLTBRootPage.String():         "HowLongToBeat Root Page",
	vangogh_local_data.LicenceProducts.String():      "Licence Products",
	vangogh_local_data.Orders.String():               "Orders",
	vangogh_local_data.PCGWEngine.String():           "PCGamingWiki Engine",
	vangogh_local_data.PCGWExternalLinks.String():    "PCGamingWiki External Links",
	vangogh_local_data.PCGWPageId.String():           "PCGamingWiki PageId",
	vangogh_local_data.SteamAppNews.String():         "Steam App News",
	vangogh_local_data.SteamReviews.String():         "Steam Reviews",
	vangogh_local_data.SteamStorePage.String():       "Steam Store Page",
	vangogh_local_data.UserWishlistProducts.String(): "User Wishlist Products",

	vangogh_local_data.MacOS.String():   "macOS",
	vangogh_local_data.Linux.String():   "Linux",
	vangogh_local_data.Windows.String(): "Windows",
}

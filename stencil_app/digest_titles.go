package stencil_app

import "github.com/arelate/vangogh_local_data"

var DigestTitles = map[string]string{
	//vangogh_local_data.OperatingSystem
	"macos":   "macOS",
	"linux":   "Linux",
	"windows": "Windows",
	// vangogh_local_data.OwnedProperty, vangogh_local_data.WishlistedProperty, ...
	vangogh_local_data.TrueValue:  "Yes",
	vangogh_local_data.FalseValue: "No",
	//vangogh_local_data.SortProperty
	vangogh_local_data.GlobalReleaseDateProperty:  PropertyTitles[vangogh_local_data.GlobalReleaseDateProperty],
	vangogh_local_data.GOGReleaseDateProperty:     PropertyTitles[vangogh_local_data.GOGReleaseDateProperty],
	vangogh_local_data.GOGOrderDateProperty:       PropertyTitles[vangogh_local_data.GOGOrderDateProperty],
	vangogh_local_data.TitleProperty:              PropertyTitles[vangogh_local_data.TitleProperty],
	vangogh_local_data.RatingProperty:             PropertyTitles[vangogh_local_data.RatingProperty],
	vangogh_local_data.DiscountPercentageProperty: PropertyTitles[vangogh_local_data.DiscountPercentageProperty],
	//Sort HLTBHoursToComplete
	vangogh_local_data.HLTBHoursToCompleteMainProperty: PropertyTitles[vangogh_local_data.HLTBHoursToCompleteMainProperty],
	vangogh_local_data.HLTBHoursToCompletePlusProperty: PropertyTitles[vangogh_local_data.HLTBHoursToCompletePlusProperty],
	vangogh_local_data.HLTBHoursToComplete100Property:  PropertyTitles[vangogh_local_data.HLTBHoursToComplete100Property],
	//vangogh_local_data.ProductTypeProperty
	vangogh_local_data.CatalogProducts.String():      "Catalog Products",
	vangogh_local_data.UserWishlistProducts.String(): "User Wishlist Products",
	vangogh_local_data.AccountProducts.String():      "Account Products",
	vangogh_local_data.ApiProductsV1.String():        "API Products V1",
	vangogh_local_data.ApiProductsV2.String():        "API Products V2",
	vangogh_local_data.Details.String():              "Account Product Details",
	vangogh_local_data.LicenceProducts.String():      "Licence Products",
	vangogh_local_data.Orders.String():               "Orders",
	vangogh_local_data.SteamAppNews.String():         "Steam App News",
	vangogh_local_data.SteamReviews.String():         "Steam Reviews",
	vangogh_local_data.SteamStorePage.String():       "Steam Store Page",
}

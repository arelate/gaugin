package view_models

import "github.com/arelate/vangogh_local_data"

var digestTitles = map[string]string{
	//vangogh_local_data.OperatingSystem
	"macos":   "macOS",
	"linux":   "Linux",
	"windows": "Windows",
	// vangogh_local_data.OwnedProperty, vangogh_local_data.WishlistedProperty, ...
	vangogh_local_data.TrueValue:  "Yes",
	vangogh_local_data.FalseValue: "No",
	//vangogh_local_data.SortProperty
	vangogh_local_data.GlobalReleaseDateProperty:  "Global Release Date",
	vangogh_local_data.GOGReleaseDateProperty:     "GOG.com Release Date",
	vangogh_local_data.GOGOrderDateProperty:       "GOG.com Order Date",
	vangogh_local_data.TitleProperty:              "Title",
	vangogh_local_data.RatingProperty:             "Rating",
	vangogh_local_data.DiscountPercentageProperty: "Discount Percentage",
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

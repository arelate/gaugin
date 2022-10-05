package stencil_app

import "github.com/arelate/vangogh_local_data"

var Labels = []string{
	vangogh_local_data.OwnedProperty,
	vangogh_local_data.ValidationResultProperty,
	vangogh_local_data.ProductTypeProperty,
	vangogh_local_data.TBAProperty,
	vangogh_local_data.ComingSoonProperty,
	vangogh_local_data.PreOrderProperty,
	vangogh_local_data.InDevelopmentProperty,
	vangogh_local_data.IsUsingDOSBoxProperty,
	vangogh_local_data.IsUsingScummVMProperty,
	vangogh_local_data.TagIdProperty,
	vangogh_local_data.LocalTagsProperty,
	vangogh_local_data.IsFreeProperty,
	vangogh_local_data.DiscountPercentageProperty,
	vangogh_local_data.WishlistedProperty,
}

var HiddenLabels = []string{
	vangogh_local_data.ValidationResultProperty,
}

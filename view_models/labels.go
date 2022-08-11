package view_models

import "github.com/arelate/vangogh_local_data"

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

func NewLabels(rdx map[string][]string) *labels {
	lbs := &labels{
		Wishlisted:     FlagFromRedux(rdx, vangogh_local_data.WishlistedProperty),
		Owned:          FlagFromRedux(rdx, vangogh_local_data.OwnedProperty),
		Free:           FlagFromRedux(rdx, vangogh_local_data.IsFreeProperty),
		Discounted:     FlagFromRedux(rdx, vangogh_local_data.IsDiscountedProperty),
		PreOrder:       FlagFromRedux(rdx, vangogh_local_data.PreOrderProperty),
		TBA:            FlagFromRedux(rdx, vangogh_local_data.TBAProperty),
		ComingSoon:     FlagFromRedux(rdx, vangogh_local_data.ComingSoonProperty),
		InDevelopment:  FlagFromRedux(rdx, vangogh_local_data.InDevelopmentProperty),
		IsUsingDOSBox:  FlagFromRedux(rdx, vangogh_local_data.IsUsingDOSBoxProperty),
		IsUsingScummVM: FlagFromRedux(rdx, vangogh_local_data.IsUsingScummVMProperty),
		Tags:           propertiesFromRedux(rdx, vangogh_local_data.TagIdProperty),
		LocalTags:      propertiesFromRedux(rdx, vangogh_local_data.LocalTagsProperty),
		ProductType:    propertyFromRedux(rdx, vangogh_local_data.ProductTypeProperty),
	}

	lbs.DiscountPercentage, lbs.DiscountLabel = discountPercentageLabelFromRedux(rdx)

	return lbs
}

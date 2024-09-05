package compton_data

import (
	"fmt"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/kevlar"
)

var LabelProperties = []string{
	vangogh_local_data.OwnedProperty,
	vangogh_local_data.ProductTypeProperty,
	vangogh_local_data.ComingSoonProperty,
	vangogh_local_data.PreOrderProperty,
	vangogh_local_data.InDevelopmentProperty,
	vangogh_local_data.TagIdProperty,
	vangogh_local_data.LocalTagsProperty,
	vangogh_local_data.IsFreeProperty,
	vangogh_local_data.DiscountPercentageProperty,
	vangogh_local_data.WishlistedProperty,
}

var labelTitles = map[string]string{
	vangogh_local_data.OwnedProperty:         "Own",
	vangogh_local_data.ComingSoonProperty:    "Soon",
	vangogh_local_data.PreOrderProperty:      "PO",
	vangogh_local_data.InDevelopmentProperty: "In Dev",
	vangogh_local_data.IsFreeProperty:        "Free",
	vangogh_local_data.WishlistedProperty:    "Wish",
}

func LabelTitle(id, property string, rdx kevlar.ReadableRedux) string {

	label, _ := rdx.GetLastVal(property, id)
	owned, _ := rdx.GetLastVal(vangogh_local_data.OwnedProperty, id)

	switch property {
	case vangogh_local_data.WishlistedProperty:
		fallthrough
	case vangogh_local_data.OwnedProperty:
		fallthrough
	case vangogh_local_data.PreOrderProperty:
		fallthrough
	case vangogh_local_data.ComingSoonProperty:
		fallthrough
	case vangogh_local_data.InDevelopmentProperty:
		fallthrough
	case vangogh_local_data.IsFreeProperty:
		if label == "true" {
			return labelTitles[property]
		}
		return ""
	case vangogh_local_data.ProductTypeProperty:
		if label == "GAME" {
			return ""
		}
	case vangogh_local_data.DiscountPercentageProperty:
		if owned == "true" {
			return ""
		}
		if label != "" && label != "0" {
			return fmt.Sprintf("-%s%%", label)
		}
		return ""
	case vangogh_local_data.TagIdProperty:
		if tagName, ok := rdx.GetLastVal(vangogh_local_data.TagNameProperty, label); ok {
			return tagName
		}
	case vangogh_local_data.DehydratedImageProperty:
		fallthrough
	case vangogh_local_data.DehydratedVerticalImageProperty:
		return property
	}
	return label
}

func LabelClass(id, property string, rdx kevlar.ReadableRedux) string {
	switch property {
	case vangogh_local_data.OwnedProperty:
		if res, ok := rdx.GetLastVal(vangogh_local_data.ValidationResultProperty, id); ok {
			if res == "OK" {
				return "validation-result-ok"
			} else {
				return "validation-result-err"
			}
		} else {
			return ""
		}
	}
	return ""
}

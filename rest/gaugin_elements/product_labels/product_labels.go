package product_labels

import (
	_ "embed"
	"fmt"
	"github.com/arelate/gaugin/rest/compton_data"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/compton/elements/labels"
	"github.com/boggydigital/kevlar"
)

//go:embed "style/product-labels.css"
var StyleProductLabels []byte

func FormatLabels(id string, rdx kevlar.ReadableRedux, properties ...string) []labels.FormattedLabel {
	owned := false
	if op, ok := rdx.GetLastVal(vangogh_local_data.OwnedProperty, id); ok {
		owned = op == vangogh_local_data.TrueValue
	}

	fmtLabels := make([]labels.FormattedLabel, 0, len(properties))

	for _, p := range properties {
		fmtLabels = append(fmtLabels, formatLabel(id, p, owned, rdx))
	}

	return fmtLabels
}

func formatLabel(id, property string, owned bool, rdx kevlar.ReadableRedux) labels.FormattedLabel {

	fmtLabel := labels.FormattedLabel{
		Property: property,
	}

	fmtLabel.Title, _ = rdx.GetLastVal(property, id)
	switch property {
	case vangogh_local_data.OwnedProperty:
		if res, ok := rdx.GetLastVal(vangogh_local_data.ValidationResultProperty, id); ok {
			if res == "OK" {
				fmtLabel.Class = "validation-result-ok"
			} else {
				fmtLabel.Class = "validation-result-err"
			}
		}
		fallthrough
	case vangogh_local_data.WishlistedProperty:
		fallthrough
	case vangogh_local_data.PreOrderProperty:
		fallthrough
	case vangogh_local_data.ComingSoonProperty:
		fallthrough
	case vangogh_local_data.InDevelopmentProperty:
		fallthrough
	case vangogh_local_data.IsFreeProperty:
		if fmtLabel.Title == "true" {
			fmtLabel.Title = compton_data.LabelTitles[property]
			break
		}
		fmtLabel.Title = ""
	case vangogh_local_data.ProductTypeProperty:
		if fmtLabel.Title == "GAME" {
			fmtLabel.Title = ""
			break
		}
	case vangogh_local_data.DiscountPercentageProperty:
		if owned {
			fmtLabel.Title = ""
			break
		}
		if fmtLabel.Title != "" && fmtLabel.Title != "0" {
			fmtLabel.Title = fmt.Sprintf("-%s%%", fmtLabel.Title)
			break
		}
		fmtLabel.Title = ""
	case vangogh_local_data.TagIdProperty:
		if tagName, ok := rdx.GetLastVal(vangogh_local_data.TagNameProperty, fmtLabel.Title); ok {
			fmtLabel.Title = tagName
			break
		}
	case vangogh_local_data.DehydratedImageProperty:
		fallthrough
	case vangogh_local_data.DehydratedVerticalImageProperty:
		fmtLabel.Title = property
	}
	return fmtLabel
}

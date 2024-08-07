package stencil_app

import (
	"fmt"
	"github.com/arelate/gaugin/data"
	"github.com/arelate/gaugin/paths"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/kevlar"
	"strconv"
	"strings"
)

var labelTitles = map[string]string{
	vangogh_local_data.OwnedProperty:         "Own",
	vangogh_local_data.ComingSoonProperty:    "Soon",
	vangogh_local_data.PreOrderProperty:      "PO",
	vangogh_local_data.InDevelopmentProperty: "In Dev",
	vangogh_local_data.IsFreeProperty:        "Free",
	vangogh_local_data.WishlistedProperty:    "Wish",
}

func ownedValidationResult(id string, rdx kevlar.ReadableRedux) (string, bool) {
	return rdx.GetLastVal(vangogh_local_data.ValidationResultProperty, id)
}

func ReviewClass(sr string) string {
	if strings.Contains(sr, "Positive") {
		return "positive"
	} else if strings.Contains(sr, "Negative") {
		return "negative"
	} else {
		return "neutral"
	}
}

func fmtAction(id, property, link string, rdx kevlar.ReadableRedux) string {

	owned, _ := rdx.GetLastVal(vangogh_local_data.OwnedProperty, id)

	switch property {
	case vangogh_local_data.WishlistedProperty:
		if owned == "true" {
			return ""
		}
		switch link {
		case "true":
			return "Remove"
		case "false":
			return "Add"
		}
	case vangogh_local_data.TagIdProperty:
		if owned != "true" {
			return ""
		}
		return "Edit"
	case vangogh_local_data.LocalTagsProperty:
		return "Edit"
	}
	return ""
}

func fmtActionHref(id, property, link string, rdx kevlar.ReadableRedux) string {

	owned, _ := rdx.GetLastVal(vangogh_local_data.OwnedProperty, id)

	switch property {
	case vangogh_local_data.WishlistedProperty:
		switch link {
		case "Add":
			return "/wishlist/add?id=" + id
		case "Remove":
			return "/wishlist/remove?id=" + id
		}
	case vangogh_local_data.TagIdProperty:
		if owned != "true" {
			return ""
		}
		return "/tags/edit?id=" + id
	case vangogh_local_data.LocalTagsProperty:
		return "/local-tags/edit?id=" + id
	}
	return ""
}

func fmtClass(id, property, link string, rdx kevlar.ReadableRedux) string {
	switch property {
	case vangogh_local_data.OwnedProperty:
		if res, ok := ownedValidationResult(id, rdx); ok {
			if res == "OK" {
				return "validation-result-ok"
			} else {
				return "validation-result-err"
			}
		} else {
			return ""
		}
	case vangogh_local_data.SteamReviewScoreDescProperty:
		return ReviewClass(link)
	case vangogh_local_data.RatingProperty:
		return ReviewClass(fmtGOGRating(link))
	case vangogh_local_data.HLTBReviewScoreProperty:
		return ReviewClass(fmtHLTBRating(link))
	case vangogh_local_data.SteamDeckAppCompatibilityCategoryProperty:
		return link
	}
	return ""
}

func fmtHref(_, property, link string, _ kevlar.ReadableRedux) string {
	switch property {
	case vangogh_local_data.GOGOrderDateProperty:
		link = justTheDate(link)
	case vangogh_local_data.PublishersProperty:
		fallthrough
	case vangogh_local_data.DevelopersProperty:
		return fmt.Sprintf("/search?%s=%s&sort=global-release-date&desc=true", property, link)
	case vangogh_local_data.IncludesGamesProperty:
		fallthrough
	case vangogh_local_data.IsIncludedByGamesProperty:
		fallthrough
	case vangogh_local_data.RequiresGamesProperty:
		fallthrough
	case vangogh_local_data.IsRequiredByGamesProperty:
		return paths.ProductId(link)
	case vangogh_local_data.RatingProperty:
		return ""
	case vangogh_local_data.DiscountPercentageProperty:
		return ""
	case vangogh_local_data.PriceProperty:
		return ""
	case data.GauginGOGLinksProperty:
		fallthrough
	case data.GauginOtherLinksProperty:
		fallthrough
	case data.GauginSteamLinksProperty:
		if _, pv, ok := strings.Cut(link, "="); ok {
			return pv
		}
	case vangogh_local_data.HLTBHoursToCompleteMainProperty:
		fallthrough
	case vangogh_local_data.HLTBHoursToCompletePlusProperty:
		fallthrough
	case vangogh_local_data.HLTBHoursToComplete100Property:
		return ""
	case vangogh_local_data.HLTBReviewScoreProperty:
		return ""
	case vangogh_local_data.EnginesBuildsProperty:
		return ""
	case vangogh_local_data.DehydratedImageProperty:
		fallthrough
	case vangogh_local_data.DehydratedVerticalImageProperty:
		return link
	}
	return fmt.Sprintf("/search?%s=%s", property, link)
}

func justTheDate(s string) string {
	return strings.Split(s, " ")[0]
}

func fmtLabel(id, property, link string, rdx kevlar.ReadableRedux) string {

	label := link
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
		if link == "true" {
			return labelTitles[property]
		}
		return ""
	case vangogh_local_data.ProductTypeProperty:
		if link == "GAME" {
			return ""
		}
	case vangogh_local_data.DiscountPercentageProperty:
		if owned == "true" {
			return ""
		}
		if link != "" && link != "0" {
			return fmt.Sprintf("-%s%%", link)
		}
		return ""
	case vangogh_local_data.TagIdProperty:
		title, ok := rdx.GetLastVal(vangogh_local_data.TagNameProperty, link)
		if !ok {
			title = link
		}
		return title
	case vangogh_local_data.DehydratedImageProperty:
		fallthrough
	case vangogh_local_data.DehydratedVerticalImageProperty:
		return property
	}
	return label
}

func fmtTitle(id, property, link string, rdx kevlar.ReadableRedux) string {
	title := link

	owned, _ := rdx.GetLastVal(vangogh_local_data.OwnedProperty, id)
	isFree, _ := rdx.GetLastVal(vangogh_local_data.IsFreeProperty, id)

	switch property {
	case vangogh_local_data.WishlistedProperty:
		if owned == "true" {
			return ""
		}
		if link == vangogh_local_data.TrueValue {
			return "Yes"
		} else {
			return "No"
		}
	case vangogh_local_data.IncludesGamesProperty:
		fallthrough
	case vangogh_local_data.IsIncludedByGamesProperty:
		fallthrough
	case vangogh_local_data.RequiresGamesProperty:
		fallthrough
	case vangogh_local_data.IsRequiredByGamesProperty:
		var ok bool
		title, ok = rdx.GetLastVal(vangogh_local_data.TitleProperty, link)
		if !ok {
			title = link
		}
	case vangogh_local_data.GOGOrderDateProperty:
		title = justTheDate(link)
	case vangogh_local_data.LanguageCodeProperty:
		title = fmt.Sprintf("%s %s", LanguageCodeFlag(link), LanguageCodeTitle(link))
	case vangogh_local_data.RatingProperty:
		title = fmtGOGRating(link)
	case vangogh_local_data.TagIdProperty:
		var ok bool
		title, ok = rdx.GetLastVal(vangogh_local_data.TagNameProperty, link)
		if !ok {
			title = link
		}
	case vangogh_local_data.PriceProperty:
		if isFree == "true" {
			return ""
		}
	case vangogh_local_data.HLTBHoursToCompleteMainProperty:
		fallthrough
	case vangogh_local_data.HLTBHoursToCompletePlusProperty:
		fallthrough
	case vangogh_local_data.HLTBHoursToComplete100Property:
		return strings.TrimLeft(link, "0") + " hrs"
	case data.GauginGOGLinksProperty:
		fallthrough
	case data.GauginOtherLinksProperty:
		fallthrough
	case data.GauginSteamLinksProperty:
		if pt, _, ok := strings.Cut(link, "="); ok {
			title = PropertyTitles[pt]
		}
	case vangogh_local_data.HLTBReviewScoreProperty:
		if link == "0" {
			return ""
		}
		return fmtHLTBRating(link)

	}

	return title
}

func fmtGOGRating(rs string) string {
	rd := ""
	if ri, err := strconv.ParseInt(rs, 10, 32); err == nil {
		rd = ratingDesc(ri * 2)
		if ri > 0 {
			rd += fmt.Sprintf(" (%.1f)", float32(ri)/10.0)
		}
	}
	return rd
}

func fmtHLTBRating(rs string) string {
	rd := ""
	if ri, err := strconv.ParseInt(rs, 10, 32); err == nil {
		rd = ratingDesc(ri)
		if ri > 0 {
			rd += fmt.Sprintf(" (%d)", ri)
		}
	}
	return rd
}

func ratingDesc(ri int64) string {
	rd := "Not Rated"
	if ri >= 95 {
		rd = "Overwhelming Positive"
	} else if ri >= 85 {
		rd = "Very Positive"
	} else if ri >= 80 {
		rd = "Positive"
	} else if ri >= 70 {
		rd = "Mostly Positive"
	} else if ri >= 40 {
		rd = "Mixed"
	} else if ri >= 20 {
		rd = "Mostly Negative"
	} else if ri > 0 {
		rd = "Negative"
	}
	return rd
}

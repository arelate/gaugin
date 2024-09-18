package compton_fragments

import (
	"fmt"
	"github.com/arelate/gaugin/data"
	"github.com/arelate/gaugin/paths"
	"github.com/arelate/gaugin/rest/compton_data"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/align"
	"github.com/boggydigital/compton/consts/direction"
	"github.com/boggydigital/compton/elements/els"
	"github.com/boggydigital/compton/elements/flex_items"
	"github.com/boggydigital/compton/elements/grid_items"
	"github.com/boggydigital/compton/elements/svg_use"
	"github.com/boggydigital/compton/elements/title_values"
	"github.com/boggydigital/kevlar"
	"slices"
	"strconv"
	"strings"
)

func ProductProperties(r compton.Registrar, id string, rdx kevlar.ReadableRedux) compton.Element {
	grid := grid_items.GridItems(r).JustifyContent(align.Center)

	for _, property := range compton_data.ProductProperties {
		if slices.Contains(compton_data.ProductHiddenProperties, property) {
			continue
		}

		propertyTitle := compton_data.PropertyTitles[property]

		// operating systems are row of SVG icons that are added separately from other properties
		if property == vangogh_local_data.OperatingSystemsProperty {
			tv := title_values.TitleValues(r, propertyTitle)
			row := flex_items.FlexItems(r, direction.Row).JustifyContent(align.Start)
			tv.Append(row)
			if values, ok := rdx.GetAllValues(property, id); ok {
				for _, os := range vangogh_local_data.ParseManyOperatingSystems(values) {
					osLink := els.A(valueHref(property, os.String()))
					osLink.Append(svg_use.SvgUse(r, compton_data.OperatingSystemSymbols[os]))
					row.Append(osLink)
				}
			}
			grid.Append(tv)
			continue
		}

		tv := title_values.TitleValues(r, propertyTitle)

		firstValue := ""
		hasContent := false

		if values, ok := rdx.GetAllValues(property, id); ok && len(values) > 0 {
			fmtValues := make([]string, 0, len(values))
			fmtValueLinks := make(map[string]string)
			for _, value := range values {
				fmtValue := valueTitle(id, property, value, rdx)
				fmtHref := valueHref(property, value)
				if fmtHref != "" && fmtValue != "" {
					fmtValueLinks[fmtValue] = fmtHref
				} else if fmtValue != "" {
					fmtValues = append(fmtValues, fmtValue)
				}
			}

			if len(fmtValues) > 0 {
				tv.AppendTextValues(fmtValues...)
				hasContent = true
			} else if len(fmtValueLinks) > 0 {
				if len(fmtValueLinks) < 4 {
					tv.AppendLinkValues(fmtValueLinks)
				} else {
					summaryElement := els.SpanText(fmt.Sprintf("Show all %d...", len(fmtValueLinks)))
					summaryElement.AddClass("action")
					ds := els.Details().AppendSummary(summaryElement)
					ds.AddClass("many-values")
					row := flex_items.FlexItems(r, direction.Row).JustifyContent(align.Start)
					for link, href := range fmtValueLinks {
						row.Append(els.AText(link, href))
					}
					ds.Append(row)
					tv.AppendValues(ds)
				}
				hasContent = true
			}

			if len(values) > 0 {
				firstValue = values[0]
			}
		}

		if fmtValueClass := valueClass(id, property, firstValue, rdx); fmtValueClass != "" && tv != nil {
			tv.AddClass(fmtValueClass)
		}
		if fmtValueAction := valueAction(id, property, firstValue, rdx); fmtValueAction != "" && tv != nil {
			fmtValueActionHref := valueActionHref(id, property, fmtValueAction, rdx)
			actionLink := els.AText(fmtValueAction, fmtValueActionHref)
			actionLink.AddClass("action")
			tv.AppendValues(actionLink)
			hasContent = true
		}

		if hasContent {
			grid.Append(tv)
		}

	}

	return grid
}

func valueAction(id, property, link string, rdx kevlar.ReadableRedux) string {

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

func valueActionHref(id, property, action string, rdx kevlar.ReadableRedux) string {

	owned, _ := rdx.GetLastVal(vangogh_local_data.OwnedProperty, id)

	switch property {
	case vangogh_local_data.WishlistedProperty:
		switch action {
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

func valueHref(property, link string) string {
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

func valueTitle(id, property, link string, rdx kevlar.ReadableRedux) string {
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
		title = fmt.Sprintf("%s %s", compton_data.LanguageCodeFlag(link), compton_data.LanguageCodeTitle(link))
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
			title = compton_data.PropertyTitles[pt]
		}
	case vangogh_local_data.HLTBReviewScoreProperty:
		if link == "0" {
			return ""
		}
		return fmtHLTBRating(link)
	}

	return title
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

func valueClass(id, property, link string, rdx kevlar.ReadableRedux) string {
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

func justTheDate(s string) string {
	return strings.Split(s, " ")[0]
}

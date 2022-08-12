package view_models

import (
	"fmt"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/yt_urls"
	"html/template"
	"strconv"
	"strings"
	"time"
)

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"transitiveDst":         transitiveDst,
		"transitiveSrc":         transitiveSrc,
		"formatBytes":           formatBytes,
		"formatDate":            formatDate,
		"unixDateFormat":        unixDateFormat,
		"unixDateTimeFormat":    unixDateTimeFormat,
		"ratingPercent":         ratingPercent,
		"toLower":               toLower,
		"downloadTitle":         downloadTitle,
		"hasLabel":              hasLabel,
		"showPrice":             showPrice,
		"propertyTitle":         propertyTitle,
		"hasDownloads":          hasDownloads,
		"languageCodeFlag":      languageCodeFlag,
		"youtubeLink":           youtubeLink,
		"shouldShowWishlist":    shouldShowWishlist,
		"canAddToWishlist":      canAddToWishlist,
		"canRemoveFromWishlist": canRemoveFromWishlist,
		"anchorId":              anchorId,
		"downloadTypeTitle":     downloadTypeTitle,
	}
}

const (
	transitiveOpen  = " ("
	transitiveClose = ")"
)

func downloadTitle(dl vangogh_local_data.Download) string {
	switch dl.Type {
	case vangogh_local_data.Installer:
		fallthrough
	case vangogh_local_data.Movie:
		fallthrough
	case vangogh_local_data.DLC:
		if !strings.HasPrefix(dl.Name, dl.ProductTitle) {
			return fmt.Sprintf("%s %s", dl.ProductTitle, dl.Name)
		} else {
			return dl.Name
		}
	case vangogh_local_data.Extra:
		return dl.Name
	default:
		return ""
	}
}

func toLower(s string) string {
	return strings.ToLower(s)
}

func transitiveDst(s string) string {
	dst := s
	if strings.Contains(s, transitiveOpen) {
		dst = s[:strings.LastIndex(s, transitiveOpen)]
	}
	return dst
}

func transitiveSrc(s string) string {
	src := ""
	if strings.Contains(s, transitiveOpen) {
		from, to := strings.LastIndex(s, transitiveOpen)+len(transitiveOpen), strings.Index(s, transitiveClose)
		if from > to {
			to = strings.LastIndex(s, transitiveClose)
			if from > to {
				from = 0
				to = len(s) - 1
			}
		}
		src = s[from:to]
	}
	return src
}

func formatBytes(b int) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}

func formatDate(d string) string {
	if hd, err := time.Parse("2006-01-02", d); err == nil {
		return hd.Format("January 2, 2006")
	}
	return ""
}

func unixDateFormat(d int) string {
	return time.Unix(int64(d), 0).Format("Jan 2, 2006")
}

func unixDateTimeFormat(d int64) string {
	return time.Unix(d, 0).Format(time.RFC1123)
}

func ratingPercent(r string) int {
	if r == "" {
		return 0
	}
	if v, err := strconv.Atoi(r); err != nil {
		return 0
	} else {
		return v * 2
	}
}

func hasLabel(lbs labels) bool {
	return lbs.Owned ||
		lbs.Wishlisted ||
		lbs.PreOrder ||
		lbs.ComingSoon ||
		lbs.TBA ||
		lbs.InDevelopment ||
		lbs.ProductType != "GAME" ||
		lbs.IsUsingDOSBox ||
		lbs.IsUsingScummVM ||
		lbs.Free ||
		lbs.Discounted ||
		len(lbs.Tags) > 0
}

func showPrice(pvm product) bool {
	if pvm.Labels.Free ||
		pvm.Labels.TBA ||
		pvm.Labels.Owned ||
		pvm.Price == "" {
		return false
	}
	return true
}

func propertyTitle(p string) string {
	return propertyTitles[p]
}

func hasDownloads(pd *ProductDownloads) bool {
	if pd == nil {
		return false
	}
	return len(pd.Installers) > 0 ||
		len(pd.DLCs) > 0 ||
		len(pd.Extras) > 0
}

func languageCodeFlag(lc string) string {
	if flag, ok := languageFlags[lc]; ok {
		return flag
	} else {
		return lc
	}
}

func youtubeLink(videoId string) string {
	return yt_urls.VideoUrl(videoId).String()
}

func shouldShowWishlist(pvm *product) bool {
	return !pvm.Labels.Owned
}

func canAddToWishlist(pvm *product) bool {
	return !pvm.Labels.Owned &&
		!pvm.Labels.Wishlisted
}

func canRemoveFromWishlist(pvm *product) bool {
	return pvm.Labels.Wishlisted
}

func anchorId(t string) string {
	return strings.Replace(strings.ToLower(t), " ", "-", -1)
}

func downloadTypeTitle(dt vangogh_local_data.DownloadType) string {
	return downloadTypeTitles[dt]
}

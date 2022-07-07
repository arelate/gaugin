package rest

import (
	"fmt"
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/yt_urls"
	"html/template"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func funcMap() template.FuncMap {
	return template.FuncMap{
		"transitiveDst":      transitiveDst,
		"transitiveSrc":      transitiveSrc,
		"formatBytes":        formatBytes,
		"justTheDate":        justTheDate,
		"formatDate":         formatDate,
		"unixDateFormat":     unixDateFormat,
		"ratingPercent":      ratingPercent,
		"gogLink":            gogLink,
		"toLower":            toLower,
		"dlTitle":            dlTitle,
		"hasLabel":           hasLabel,
		"hasTags":            hasTags,
		"showPrice":          showPrice,
		"searchPropertyName": searchPropertyName,
		"hasDownloads":       hasDownloads,
		"hasSteamLinks":      hasSteamLinks,
		"hasGOGLinks":        hasGOGLinks,
		"languageCodeFlag":   languageCodeFlag,
		"youtubeLink":        youtubeLink,
		"hasShortcuts":       hasShortcuts,
	}
}

const (
	transitiveOpen  = " ("
	transitiveClose = ")"
)

func dlTitle(dl vangogh_local_data.Download) string {
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

func justTheDate(s string) string {
	return strings.Split(s, " ")[0]
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

func gogLink(p string) string {
	u := url.URL{
		Scheme: gog_integration.HttpsScheme,
		Host:   gog_integration.WwwGogHost,
		Path:   p,
	}
	return u.String()
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

func hasTags(lbs labels) bool {
	return len(lbs.Tags) > 0 || len(lbs.LocalTags) > 0
}

func showPrice(pvm productViewModel) bool {
	if pvm.Labels.Free ||
		pvm.Labels.TBA ||
		pvm.Labels.Owned ||
		pvm.Price == "" {
		return false
	}
	return true
}

func searchPropertyName(p string) string {
	return searchPropertyNames[p]
}

func hasDownloads(pd *productDownloads) bool {
	if pd == nil {
		return false
	}
	return len(pd.Installers) > 0 ||
		len(pd.DLCs) > 0 ||
		len(pd.Extras) > 0
}

func hasSteamLinks(pvm *productViewModel) bool {
	return pvm.SteamAppId != ""
}

func hasGOGLinks(pvm *productViewModel) bool {
	return pvm.StoreUrl != "" ||
		pvm.ForumUrl != "" ||
		pvm.SupportUrl != ""
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

func hasShortcuts(pvm *productViewModel) bool {
	return pvm.HasDescription ||
		pvm.HasChangelog ||
		pvm.HasScreenshots ||
		pvm.HasVideos ||
		pvm.HasDownloads ||
		pvm.HasSteamAppNews
}

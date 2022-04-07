package rest

import (
	"fmt"
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_local_data"
	"html/template"
	"net/url"
	"strconv"
	"strings"
)

func funcMap() template.FuncMap {
	return template.FuncMap{
		"productTitle":  productTitle,
		"productId":     productId,
		"formatBytes":   formatBytes,
		"justTheDate":   justTheDate,
		"ratingPercent": ratingPercent,
		"gogLink":       gogLink,
		"toLower":       toLower,
		"dlTitle":       dlTitle,
	}
}

const titleIdSep = " ("

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

func productTitle(s string) string {
	if strings.Contains(s, titleIdSep) {
		return s[:strings.LastIndex(s, titleIdSep)]
	}
	return s
}

func productId(s string) string {
	if strings.Contains(s, titleIdSep) {
		return s[strings.LastIndex(s, titleIdSep)+len(titleIdSep) : len(s)-1]
	}
	return ""
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

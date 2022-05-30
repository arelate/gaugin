package rest

import (
	"fmt"
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_local_data"
	"html/template"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func funcMap() template.FuncMap {
	return template.FuncMap{
		"productTitle":       productTitle,
		"productId":          productId,
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

func hasLabel(lpvm listProductViewModel) bool {
	return lpvm.Owned ||
		lpvm.Wishlisted ||
		lpvm.PreOrder ||
		lpvm.ComingSoon ||
		lpvm.TBA ||
		lpvm.InDevelopment ||
		lpvm.ProductType != "GAME" ||
		lpvm.IsUsingDOSBox ||
		lpvm.IsUsingScummVM ||
		lpvm.Free ||
		lpvm.Discounted ||
		len(lpvm.Tags) > 0
}

func hasTags(pvm productViewModel) bool {
	return len(pvm.Tags) > 0 || len(pvm.LocalTags) > 0
}

func showPrice(pvm productViewModel) bool {
	if pvm.Free ||
		pvm.TBA ||
		pvm.Owned ||
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

package api

import (
	"fmt"
	"github.com/arelate/vangogh_local_data"
	"html/template"
	"strconv"
	"strings"
)

func funcMap() template.FuncMap {
	return template.FuncMap{
		"productTitle":  productTitle,
		"productId":     productId,
		"downloadTitle": downloadTitle,
		"formatBytes":   formatBytes,
		"justTheDate":   justTheDate,
		"ratingPercent": ratingPercent,
	}
}

const titleIdSep = " ("

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

func downloadTitle(d vangogh_local_data.Download) string {
	return d.String()
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

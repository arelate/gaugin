package view_models

import (
	"fmt"
	"github.com/arelate/gaugin/stencil_app"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/yet_urls/youtube_urls"
	"html/template"
	"strings"
	"time"
)

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"formatBytes":       formatBytes,
		"unixDateFormat":    unixDateFormat,
		"toLower":           toLower,
		"downloadTitle":     downloadTitle,
		"hasDownloads":      hasDownloads,
		"languageCodeFlag":  stencil_app.LanguageCodeFlag,
		"youtubeLink":       youtubeLink,
		"downloadTypeTitle": downloadTypeTitle,
	}
}

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

func unixDateFormat(d int64) string {
	return time.Unix(d, 0).Format("Jan 2, 2006")
}

func hasDownloads(pd *ProductDownloads) bool {
	if pd == nil {
		return false
	}
	return len(pd.Installers) > 0 ||
		len(pd.DLCs) > 0 ||
		len(pd.Extras) > 0
}

func youtubeLink(videoId string) string {
	return youtube_urls.VideoUrl(videoId).String()
}

func downloadTypeTitle(dt vangogh_local_data.DownloadType) string {
	return downloadTypeTitles[dt]
}

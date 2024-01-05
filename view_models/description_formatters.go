package view_models

import (
	"github.com/arelate/gaugin/paths"
	"github.com/arelate/vangogh_local_data"
	"net/url"
	"path"
	"strings"
)

func rewriteItemsLinks(desc string) string {

	itemsUrls := vangogh_local_data.ExtractDescItems(desc)

	for _, itemUrl := range itemsUrls {
		if u, err := url.Parse(itemUrl); err != nil {
			continue
		} else {
			ggUrl := "/items" + u.Path
			desc = strings.Replace(desc, itemUrl, ggUrl, -1)
		}
	}

	return desc
}

func rewriteGameLinks(desc string) string {
	gameLinks := vangogh_local_data.ExtractGameLinks(desc)

	for _, gameLink := range gameLinks {
		if u, err := url.Parse(gameLink); err != nil {
			continue
		} else {
			_, slug := path.Split(u.Path)
			ggUrl := paths.ProductSlug(slug)
			desc = strings.Replace(desc, gameLink, ggUrl, -1)
		}
	}

	return desc
}

func rewriteLinksAsTargetTop(desc string) string {
	return strings.Replace(desc, "<a ", "<a target='_top' ", -1)
}

func rewriteVideoAsInline(desc string) string {
	return strings.Replace(desc, "<video ", "<video playsinline ", -1)
}

func fixQuotes(desc string) string {
	return strings.Replace(desc, "”", "\"", -1)
}

func replaceDataFallbackUrls(desc string) string {
	return strings.Replace(desc, "data-fallbackurl", "poster", -1)
}

const doubleNewLineChar = "\n\n"
const newLineChar = "\n"
const emDashCode = "\u2013"

// implicitToExplicitList looks for embedded characters
// that GOG.com is using for <ul> lists creation, e.g.
// https://www.gog.com/en/game/deaths_gambit
// and replaces that segment with explicit unordered lists.
// Currently known characters are listed as consts above this func.
func implicitToExplicitList(text string) string {
	var items []string
	if strings.Contains(text, doubleNewLineChar) {
		items = strings.Split(text, doubleNewLineChar)
	} else if strings.Contains(text, newLineChar) {
		items = strings.Split(text, newLineChar)
	} else if strings.Contains(text, emDashCode) {
		items = strings.Split(text, emDashCode)
	}

	if len(items) > 0 {
		builder := strings.Builder{}
		builder.WriteString("<ul>")
		for _, item := range items {
			builder.WriteString("<li>" + item + "</li>")
		}
		builder.WriteString("</ul>")
		text = builder.String()
	}

	return text
}

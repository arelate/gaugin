package view_models

import (
	"fmt"
	"github.com/boggydigital/yt_urls"
	"regexp"
	"strings"
)

var preFormatters = map[string]string{
	"\n\n": "\n",
}

var formatters = map[string]string{
	"[<a":                   "<a",
	"/a>]":                  "/a>",
	"[b]":                   "<b>",
	"[/b]":                  "</b>",
	"[B]":                   "<b>",
	"[/B]":                  "</b>",
	"[<em>":                 "<em>",
	"</em>]":                "</em>",
	"[u]":                   "<u>",
	"[/u]":                  "</u>",
	"[i]":                   "<i>",
	"[/i]":                  "</i>",
	"[I]":                   "<i>",
	"[/I]":                  "</i>",
	"[strike]":              "<s>",
	"[/strike]":             "</s>",
	"[h1]":                  "<h1>",
	"[/h1]":                 "</h1>",
	"[H1]":                  "<h1>",
	"[/H1]":                 "</h1>",
	"[h2]":                  "<h2>",
	"[/h2]":                 "</h2>",
	"[H2]":                  "<h2>",
	"[/H2]":                 "</h2>",
	"[h3]":                  "<h3>",
	"[/h3]":                 "</h3>",
	"[H3]":                  "<h3>",
	"[/H3]":                 "</h3>",
	"[list]":                "<ul>",
	"[List]":                "<ul>",
	"[*]":                   "<li>",
	"[/*]":                  "</li>",
	"[/list]":               "</ul>",
	"[/List]":               "</ul>",
	"[olist]":               "<ol>",
	"[/olist]":              "</ol>",
	"\n":                    "<br>",
	"{STEAM_CLAN_IMAGE}":    "https://cdn.cloudflare.steamstatic.com/steamcommunity/public/images/clans/",
	"[img]":                 "<img src='",
	"[/img]":                "'/>",
	"[img=":                 "[",
	"[IMG]":                 "<img src='",
	"[/IMG]":                "'/>",
	"[hr][/hr]":             "<hr/>",
	"[spoiler]":             "<mark class='spoiler' tabindex='0' title='press and hold to reveal'>",
	"[/spoiler]":            "</mark>",
	"[quote]":               "&quot;",
	"[/quote]":              "&quot;",
	"[code]":                "<pre>",
	"[/code]":               "</pre>",
	"[url]":                 "<a href='",
	"[url=":                 "",
	"[and]":                 "&",
	"[/url]":                "'/>",
	"[table]":               "<table>",
	"[/table]":              "</table>",
	"[at]":                  "@",
	"[dot]":                 ".",
	"[redacted]":            "<mark class='spoiler' tabindex='0' title='press and hold to reveal'>",
	"[/redacted]":           "</mark>",
	"[cms-block]":           "",
	"[tr]":                  "<tr>",
	"[/tr]":                 "</tr>",
	"[td]":                  "<td>",
	"[/td]":                 "</td>",
	"[th]":                  "<th>",
	"[/th]":                 "</th>",
	"[expand type=details]": "<details><summary>(Tap to expand)</summary>",
	"[/expand]":             "</details>",
}

func steamNewsToHTML(c string) string {

	for t := range preFormatters {
		c = strings.Replace(c, t, preFormatters[t], -1)
	}

	for t := range formatters {
		c = strings.Replace(c, t, formatters[t], -1)
	}

	c = replaceUrls(c)
	c = replaceYoutubePreviews(c)

	return c
}

var (
	steamNewsItemUrl    = regexp.MustCompile(`(?s)\[url=.*?\[/url]`)
	steamYoutubePreview = regexp.MustCompile(`(?s)\[previewyoutube=.*?\[/previewyoutube]`)
)

func extractUrls(c string) []string {
	return steamNewsItemUrl.FindAllString(c, -1)
}

func extractYoutubePreviews(c string) []string {
	return steamYoutubePreview.FindAllString(c, -1)
}

func replaceUrls(c string) string {
	urls := extractUrls(c)
	for _, u := range urls {
		nu := strings.Replace(u, "[url=", "<a target='_top' href='", 1)
		nu = strings.Replace(nu, "]", "'>", 1)
		nu = strings.Replace(nu, "[/url]", "</a>", 1)
		c = strings.Replace(c, u, nu, 1)
	}
	return c
}

func replaceYoutubePreviews(c string) string {
	youtubePreviews := extractYoutubePreviews(c)
	for _, yp := range youtubePreviews {
		vidId := strings.TrimPrefix(yp, "[previewyoutube=")
		vidId = strings.TrimSuffix(vidId, "][/previewyoutube]")
		vidId = strings.Split(vidId, ";")[0]

		nyp := fmt.Sprintf("<p><a target='_top' href='%s'>YouTube Link</a></p>", yt_urls.VideoUrl(vidId))

		c = strings.Replace(c, yp, nyp, -1)
	}
	return c
}

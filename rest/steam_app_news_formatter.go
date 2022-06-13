package rest

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
	"[b]":                "<b>",
	"[/b]":               "</b>",
	"[u]":                "<u>",
	"[/u]":               "</u>",
	"[i]":                "<i>",
	"[/i]":               "</i>",
	"[h1]":               "<h1>",
	"[/h1]":              "</h1>",
	"[h2]":               "<h2>",
	"[/h2]":              "</h2>",
	"[h3]":               "<h3>",
	"[/h3]":              "</h3>",
	"[list]":             "<ul>",
	"[*]":                "<li>",
	"[/list]":            "</ul>",
	"\n":                 "<br>",
	"{STEAM_CLAN_IMAGE}": "https://cdn.cloudflare.steamstatic.com/steamcommunity/public/images/clans/",
	"[img]":              "<img src='",
	"[/img]":             "'/>",
	"[hr][/hr]":          "<hr/>",
}

func steamAppNewsToHTML(c string) string {

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

package rest

import (
	"regexp"
	"strings"
)

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
}

func steamAppNewsToHTML(c string) string {

	urls := extractUrls(c)
	for _, u := range urls {
		nu := strings.Replace(u, "[url=", "<a target='_top' href='", 1)
		nu = strings.Replace(nu, "]", "'>", 1)
		nu = strings.Replace(nu, "[/url]", "</a>", 1)
		c = strings.Replace(c, u, nu, 1)
	}

	for t := range formatters {
		c = strings.Replace(c, t, formatters[t], -1)
	}

	return c
}

var steamNewsItemUrl = regexp.MustCompile(`\[url=[\w:/*-. ]*][\w:\/\[\]*-. !?,;]*?\[/url]`)

func extractUrls(c string) []string {
	return steamNewsItemUrl.FindAllString(c, -1)
}

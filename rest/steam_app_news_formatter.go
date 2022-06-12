package rest

import "strings"

var formatters = map[string]string{
	"[b]":     "<b>",
	"[/b]":    "</b>",
	"[u]":     "<u>",
	"[/u]":    "</u>",
	"[h1]":    "<h1>",
	"[/h1]":   "</h1>",
	"[h2]":    "<h2>",
	"[/h2]":   "</h2>",
	"[h3]":    "<h3>",
	"[/h3]":   "</h3>",
	"[list]":  "<ul>",
	"[*]":     "<li>",
	"[/list]": "</ul>",
	"\n":      "<br>",
	//since CSP will not allow images from another host - converting images to links for now
	"[img]":              "<a target='_top' href='",
	"[/img]":             "'>External image</a>",
	"{STEAM_CLAN_IMAGE}": "https://cdn.cloudflare.steamstatic.com/steamcommunity/public/images/clans/",
}

func steamAppNewsToHTML(c string) string {
	for t := range formatters {
		c = strings.Replace(c, t, formatters[t], -1)
	}

	return c
}

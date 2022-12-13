package view_models

import (
	"github.com/arelate/southern_light/steam_integration"
	"html/template"
	"strings"
)

type steamNews struct {
	Context   string
	Count     uint32
	NewsItems []steamNewsItem
}

type steamNewsItem struct {
	Title     string
	Date      int64
	Author    string
	Url       string
	Tags      string
	FeedLabel string
	Contents  template.HTML
}

func NewSteamNews(san *steam_integration.AppNews) *steamNews {
	sanvm := &steamNews{
		Context:   "iframe",
		Count:     san.Count,
		NewsItems: make([]steamNewsItem, 0, len(san.NewsItems)),
	}

	for _, ni := range san.NewsItems {
		sanvm.NewsItems = append(sanvm.NewsItems, steamNewsItem{
			Title:     ni.Title,
			Date:      ni.Date,
			Author:    ni.Author,
			Url:       ni.Url,
			Tags:      strings.Join(ni.Tags, ","),
			FeedLabel: ni.FeedLabel,
			Contents:  template.HTML(steamNewsToHTML(ni.Contents)),
		})
	}

	return sanvm
}

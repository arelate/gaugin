package view_models

import (
	"github.com/arelate/steam_integration"
	"html/template"
)

type steamReviews struct {
	Context string
	Count   int
	Reviews []steamReview
}

type steamReview struct {
	Author                   string
	Language                 string
	Created, Updated         string
	VotedUp                  bool
	VotesUp, VotesFunny      int
	SteamPurchase            bool
	ReceivedForFree          bool
	WrittenDuringEarlyAccess bool
	Review                   template.HTML
}

func NewSteamReviews(sar *steam_integration.AppReviews) *steamReviews {
	srvm := &steamReviews{
		Context: "iframe",
		Count:   len(sar.Reviews),
		Reviews: make([]steamReview, 0, len(sar.Reviews)),
	}

	for _, rev := range sar.Reviews {
		srvm.Reviews = append(srvm.Reviews, steamReview{
			Author:                   rev.Author.SteamId,
			Language:                 rev.Language,
			Created:                  unixDateFormat(rev.TimestampCreated),
			Updated:                  unixDateFormat(rev.TimestampUpdated),
			VotedUp:                  rev.VotedUp,
			VotesUp:                  rev.VotesUp,
			VotesFunny:               rev.VotesFunny,
			SteamPurchase:            rev.SteamPurchase,
			ReceivedForFree:          rev.ReceivedForFree,
			WrittenDuringEarlyAccess: rev.WrittenDuringEarlyAccess,
			Review:                   template.HTML(rev.Review),
		})
	}

	return srvm
}

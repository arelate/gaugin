package view_models

import "github.com/arelate/southern_light/steam_integration"

type steamDeck struct {
	Category           string
	BlogUrl            string
	ResultsDisplayType []int
	Results            []string
}

func NewSteamDeck(dacr *steam_integration.DeckAppCompatibilityReport) *steamDeck {
	return &steamDeck{
		Category:           dacr.String(),
		BlogUrl:            dacr.GetSteamDeckBlogUrl(),
		ResultsDisplayType: dacr.GetResultsDisplayTypes(),
		Results:            dacr.GetResults(),
	}

}

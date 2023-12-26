package view_models

import (
	"github.com/arelate/southern_light/steam_integration"
	"html/template"
)

var messageByCategory = map[string]template.HTML{
	"Verified": "Valve’s testing indicates that this game is <span class='category verified'>Verified</span> on Steam Deck. " +
		"This game is fully functional on Steam Deck, and works great with the built-in controls and display.",
	"Playable": "Valve’s testing indicates that this game is <span class='category playable'>Playable</span> on Steam Deck. " +
		"This game is functional on Steam Deck, but might require extra effort to interact with or configure.",
	"Unsupported": "Valve’s testing indicates that 8-Bit Hordes is <span class='category unsupported'>Unsupported</span> on Steam Deck. " +
		"Some or all of this game currently doesn't function on Steam Deck.",
	"Unknown": "Unknown",
}

type steamDeck struct {
	Message            template.HTML
	BlogUrl            string
	ResultsDisplayType []int
	Results            []string
}

func NewSteamDeck(dacr *steam_integration.DeckAppCompatibilityReport) *steamDeck {
	return &steamDeck{
		Message:            messageByCategory[dacr.String()],
		BlogUrl:            dacr.GetSteamDeckBlogUrl(),
		ResultsDisplayType: dacr.GetResultsDisplayTypes(),
		Results:            dacr.GetResults(),
	}

}

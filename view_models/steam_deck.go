package view_models

import (
	"fmt"
	"github.com/arelate/southern_light/steam_integration"
	"html/template"
)

var messageByCategory = map[string]string{
	"Verified": "Valve’s testing indicates that <span class='title'>%s</span> is " +
		"<span class='category verified'>Verified</span> on Steam Deck. " +
		"This game is fully functional on Steam Deck, and works great with the built-in controls and display.",
	"Playable": "Valve’s testing indicates that <span class='title'>%s</span> is " +
		"<span class='category playable'>Playable</span> on Steam Deck. " +
		"This game is functional on Steam Deck, but might require extra effort to interact with or configure.",
	"Unsupported": "Valve’s testing indicates that <span class='title'>%s</span> is " +
		"<span class='category unsupported'>Unsupported</span> on Steam Deck. " +
		"Some or all of this game currently doesn't function on Steam Deck.",
	"Unknown": "Valve is still learning about <span class='title'>%s</span>. " +
		"We do not currently have further information regarding Steam Deck compatibility.",
}

type steamDeck struct {
	Message      template.HTML
	BlogUrl      string
	DisplayTypes []string
	Results      []string
}

func NewSteamDeck(title string, dacr *steam_integration.DeckAppCompatibilityReport) *steamDeck {

	message := fmt.Sprintf(messageByCategory[dacr.String()], title)

	return &steamDeck{
		Message:      template.HTML(message),
		BlogUrl:      dacr.GetSteamDeckBlogUrl(),
		DisplayTypes: dacr.GetResultsDisplayTypes(),
		Results:      dacr.GetResults(),
	}

}

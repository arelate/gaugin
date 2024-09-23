package rest

import (
	"fmt"
	"github.com/arelate/gaugin/rest/compton_data"
	"github.com/arelate/gaugin/rest/gaugin_styles"
	"github.com/arelate/southern_light/steam_integration"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/compton/consts/direction"
	"github.com/boggydigital/compton/elements/els"
	"github.com/boggydigital/compton/elements/flex_items"
	"github.com/boggydigital/compton/elements/iframe_expand"
	"github.com/boggydigital/nod"
	"net/http"
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

func GetSteamDeck(w http.ResponseWriter, r *http.Request) {

	// GET /steam-deck?id

	id := r.URL.Query().Get("id")

	dacr, err := getSteamDeckReport(http.DefaultClient, id)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	titleResult, err := getRedux(http.DefaultClient, id, false, vangogh_local_data.TitleProperty)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	title := ""
	if tt, ok := titleResult[id][vangogh_local_data.TitleProperty]; ok && len(tt) > 0 {
		title = tt[0]
	}

	message := fmt.Sprintf(messageByCategory[dacr.String()], title)

	section := compton_data.SteamDeckSection
	ifc := iframe_expand.IframeExpandContent(section, compton_data.SectionTitles[section]).
		AppendStyle(gaugin_styles.SteamDeckStyle)

	pageStack := flex_items.FlexItems(ifc, direction.Column)
	ifc.Append(pageStack)

	divMessage := els.DivText(message)
	divMessage.AddClass("message")
	pageStack.Append(divMessage)

	results := dacr.GetResults()

	if len(results) > 0 {
		pageStack.Append(els.Hr())
	}

	if blogUrl := dacr.GetSteamDeckBlogUrl(); blogUrl != "" {
		pageStack.Append(els.AText("Read more in the Steam blog", blogUrl))
	}

	ul := els.Ul()
	if displayTypes := dacr.GetResultsDisplayTypes(); len(displayTypes) == len(results) {
		for ii, result := range results {
			decodedResult := steam_integration.DecodeLocToken(result)
			li := els.ListItemText(decodedResult)
			li.AddClass(displayTypes[ii])
			ul.Append(li)
		}
	}
	pageStack.Append(ul)

	if err := ifc.WriteContent(w); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
	}
}

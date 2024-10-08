package compton_data

import (
	"github.com/boggydigital/compton/elements/svg_use"
)

const AppTitle = "gaugin"

const AppFavIconEmoji = "🪸"

const (
	AppNavUpdates = "Updates"
	AppNavSearch  = "Search"
)

var AppNavOrder = []string{AppNavUpdates, AppNavSearch}

var AppNavIcons = map[string]svg_use.Symbol{
	AppNavUpdates: svg_use.Sparkle,
	AppNavSearch:  svg_use.Search,
}

var AppNavLinks = map[string]string{
	AppNavUpdates: UpdatesPath,
	AppNavSearch:  SearchPath,
}

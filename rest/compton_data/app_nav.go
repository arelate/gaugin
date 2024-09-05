package compton_data

import (
	"github.com/boggydigital/compton/elements/svg_inline"
)

const (
	AppNavUpdates = "Updates"
	AppNavSearch  = "Search"
)

var AppNavOrder = []string{AppNavUpdates, AppNavSearch}

var AppNavIcons = map[string]svg_inline.Symbol{
	AppNavUpdates: svg_inline.Sparkle,
	AppNavSearch:  svg_inline.Search,
}

var AppNavLinks = map[string]string{
	AppNavUpdates: UpdatesPath,
	AppNavSearch:  SearchPath,
}

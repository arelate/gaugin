package stencil_app

import "github.com/boggydigital/stencil"

const (
	NavUpdates = "Updates"
	NavSearch  = "Search"
)

var NavItems = []string{NavUpdates, NavSearch}

var NavIcons = map[string]string{
	NavUpdates: stencil.IconSparkle,
	NavSearch:  stencil.IconSearch,
}

var NavHrefs = map[string]string{
	NavUpdates: UpdatesPath,
	NavSearch:  SearchPath,
}

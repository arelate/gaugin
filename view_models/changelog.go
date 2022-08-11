package view_models

import (
	"github.com/arelate/vangogh_local_data"
	"html/template"
)

type changelog struct {
	Context   string
	Changelog template.HTML
}

func NewChangelog(rdx map[string][]string) *changelog {
	cvm := &changelog{Context: "iframe"}

	clog := propertyFromRedux(rdx, vangogh_local_data.ChangelogProperty)
	clog = rewriteLinksAsTargetTop(clog)

	cvm.Changelog = template.HTML(clog)

	return cvm
}

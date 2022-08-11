package view_models

import "github.com/arelate/vangogh_local_data"

type tagsEdit struct {
	Context      string
	Id           string
	Title        string
	Owned        bool
	AllTags      []string
	AllLocalTags []string
	Tags         map[string]bool
	LocalTags    map[string]bool
}

func NewTagsEdit(id string, rdx map[string][]string, digests map[string][]string) *tagsEdit {
	tevm := &tagsEdit{
		Context:      "tags",
		Id:           id,
		Title:        propertyFromRedux(rdx, vangogh_local_data.TitleProperty),
		Owned:        FlagFromRedux(rdx, vangogh_local_data.OwnedProperty),
		AllTags:      propertiesFromRedux(digests, vangogh_local_data.TagIdProperty),
		AllLocalTags: propertiesFromRedux(digests, vangogh_local_data.LocalTagsProperty),
	}

	selectedTags := make(map[string]bool)
	for _, t := range propertiesFromRedux(rdx, vangogh_local_data.TagIdProperty) {
		selectedTags[t] = true
	}

	selectedLocalTags := make(map[string]bool)
	for _, t := range propertiesFromRedux(rdx, vangogh_local_data.LocalTagsProperty) {
		selectedLocalTags[t] = true
	}

	tevm.Tags = selectedTags
	tevm.LocalTags = selectedLocalTags

	return tevm
}

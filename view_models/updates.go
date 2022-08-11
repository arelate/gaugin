package view_models

import "sort"

type updates struct {
	Context         string
	Sections        []string
	SectionProducts map[string]*list
}

func NewUpdates(
	summary map[string][]string,
	rdx map[string]map[string][]string) *updates {

	sections := make([]string, 0, len(summary))
	sectionProducts := make(map[string]*list)
	for s, ids := range summary {
		sections = append(sections, s)
		sectionProducts[s] = NewListViewModel(ids, rdx)
	}

	sort.Strings(sections)

	uvm := &updates{
		Context:         "updates",
		Sections:        sections,
		SectionProducts: sectionProducts,
	}

	return uvm
}

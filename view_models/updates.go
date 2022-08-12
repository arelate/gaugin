package view_models

import (
	"github.com/arelate/vangogh_local_data"
	"sort"
	"strconv"
)

type updates struct {
	Context         string
	Sections        []string
	SectionProducts map[string]*list
	SyncCompleted   int64
}

func NewUpdates(
	summary map[string][]string,
	dataRdx map[string]map[string][]string,
	syncRdx map[string][]string) *updates {

	sections := make([]string, 0, len(summary))
	sectionProducts := make(map[string]*list)
	for s, ids := range summary {
		sections = append(sections, s)
		sectionProducts[s] = NewListViewModel(ids, dataRdx)
	}

	sort.Strings(sections)

	uvm := &updates{
		Context:         "updates",
		Sections:        sections,
		SectionProducts: sectionProducts,
	}

	scs := propertyFromRedux(syncRdx, vangogh_local_data.SyncEventsProperty)
	if syncCompleted, err := strconv.ParseInt(scs, 10, 64); err == nil {
		uvm.SyncCompleted = syncCompleted
	}

	return uvm
}

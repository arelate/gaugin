package rest

type IdPropertyValues = map[string]map[string][]string

func MergeIdPropertyValues(first, second IdPropertyValues) IdPropertyValues {
	for id, pv := range second {
		if first[id] == nil {
			first[id] = make(map[string][]string)
		}
		for p, v := range pv {
			first[id][p] = append(first[id][p], v...)
		}
	}
	return first
}

func NewIdPropertyValues(properties []string) IdPropertyValues {
	idPropertyValues := make(IdPropertyValues)
	idPropertyValues[""] = make(map[string][]string)
	for _, p := range properties {
		idPropertyValues[""][p] = nil
	}
	return idPropertyValues
}

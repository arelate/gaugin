package view_models

import (
	"github.com/arelate/vangogh_local_data"
)

func propertyFromRedux(redux map[string][]string, property string) string {
	properties := propertiesFromRedux(redux, property)
	if len(properties) > 0 {
		return properties[0]
	}
	return ""
}

func FlagFromRedux(redux map[string][]string, property string) bool {
	return propertyFromRedux(redux, property) == vangogh_local_data.TrueValue
}

func propertiesFromRedux(redux map[string][]string, property string) []string {
	if val, ok := redux[property]; ok {
		return val
	} else {
		return []string{}
	}
}

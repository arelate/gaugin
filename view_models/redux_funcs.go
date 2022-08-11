package view_models

import (
	"github.com/arelate/vangogh_local_data"
	"strconv"
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

func discountPercentageLabelFromRedux(redux map[string][]string) (int, string) {
	dp, dl := 0, ""
	dpa := propertyFromRedux(redux, vangogh_local_data.DiscountPercentageProperty)
	if dpi, err := strconv.Atoi(dpa); err == nil {
		dp = dpi
		if dp >= 80 {
			dl = "\u2158" // 4/5
		} else if dp >= 75 {
			dl = "\u00be" // 3/4
		} else if dp >= 66 {
			dl = "\u2154" // 2/3
		} else if dp >= 60 {
			dl = "\u2157" // 3/5
		} else if dp >= 50 {
			dl = "\u00bd" // 1/2
		} else if dp >= 40 {
			dl = "\u2156" // 2/5
		} else if dp >= 33 {
			dl = "\u2153" // 1/3
		} else if dp >= 25 {
			dl = "\u00bc" // 1/4
		} else if dp >= 20 {
			dl = "\u2155" // 1/5
		}
	}
	return dp, dl
}

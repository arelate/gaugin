package view_models

import "github.com/arelate/vangogh_local_data"

type listProduct struct {
	Id               string
	Title            string
	Labels           *labels
	OperatingSystems []string
	Properties       map[string][]string
	PropertyOrder    []string
	PropertyTitles   map[string]string
}

type list struct {
	Context  string
	Products []listProduct
}

func NewListViewModel(ids []string, redux map[string]map[string][]string) *list {
	lvm := &list{
		Products: make([]listProduct, 0, len(ids)),
	}
	for _, id := range ids {
		rdx, ok := redux[id]
		if !ok {
			continue
		}
		lpvm := listProduct{
			Id:               id,
			Title:            propertyFromRedux(rdx, vangogh_local_data.TitleProperty),
			Labels:           NewLabels(rdx),
			OperatingSystems: propertiesFromRedux(rdx, vangogh_local_data.OperatingSystemsProperty),
			Properties:       make(map[string][]string),
			PropertyTitles:   propertyTitles,
			PropertyOrder:    listPropertyOrder,
		}

		for _, p := range listPropertyOrder {
			lpvm.Properties[p] = propertiesFromRedux(rdx, p)
		}

		lvm.Products = append(lvm.Products, lpvm)
	}
	return lvm
}

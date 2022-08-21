package view_models

import (
	"github.com/arelate/vangogh_local_data"
	"golang.org/x/exp/maps"
	"strconv"
	"strings"
)

type Downloads struct {
	ValidationSuccess   bool
	ValidationResults   []string
	ValidationCompleted string
	Context             string
	CurrentOS           *ProductDownloads
	OtherOS             *ProductDownloads
	Extras              *ProductDownloads
}

type ProductDownloads struct {
	CurrentOS        bool
	OperatingSystems string
	Installers       vangogh_local_data.DownloadsList
	DLCs             vangogh_local_data.DownloadsList
	Extras           vangogh_local_data.DownloadsList
}

func NewDownloads(rdx map[string][]string, clientOS vangogh_local_data.OperatingSystem, dls vangogh_local_data.DownloadsList) *Downloads {

	dvm := &Downloads{
		Context:           "iframe",
		ValidationSuccess: propertyFromRedux(rdx, vangogh_local_data.ValidationResultProperty) == vangogh_local_data.OKValue,
		ValidationResults: propertiesFromRedux(rdx, vangogh_local_data.ValidationResultProperty),
		CurrentOS: &ProductDownloads{
			OperatingSystems: clientOS.String(),
			CurrentOS:        true,
			Installers:       make(vangogh_local_data.DownloadsList, 0, len(dls)),
			DLCs:             make(vangogh_local_data.DownloadsList, 0, len(dls)),
		},
		OtherOS: &ProductDownloads{
			CurrentOS:  false,
			Installers: make(vangogh_local_data.DownloadsList, 0, len(dls)),
			DLCs:       make(vangogh_local_data.DownloadsList, 0, len(dls)),
		},
		Extras: &ProductDownloads{
			CurrentOS: false,
			Extras:    make(vangogh_local_data.DownloadsList, 0, len(dls)),
		},
	}

	ustr := propertyFromRedux(rdx, vangogh_local_data.ValidationCompletedProperty)
	if unix, err := strconv.ParseInt(ustr, 10, 64); err == nil {
		dvm.ValidationCompleted = unixDateFormat(unix)
	}

	otherOS := make(map[string]interface{})

	var osd *ProductDownloads
	for _, dl := range dls {
		if dl.OS == clientOS {
			osd = dvm.CurrentOS
		} else if dl.OS == vangogh_local_data.AnyOperatingSystem {
			osd = dvm.Extras
		} else {
			otherOS[dl.OS.String()] = nil
			osd = dvm.OtherOS
		}

		switch dl.Type {
		case vangogh_local_data.Installer:
			osd.Installers = append(osd.Installers, dl)
		case vangogh_local_data.DLC:
			osd.DLCs = append(osd.DLCs, dl)
		case vangogh_local_data.Extra:
			fallthrough
		default:
			osd.Extras = append(osd.Extras, dl)
		}
	}

	dvm.OtherOS.OperatingSystems = strings.Join(maps.Keys(otherOS), ", ")
	dvm.Extras.OperatingSystems = "Other"

	return dvm
}

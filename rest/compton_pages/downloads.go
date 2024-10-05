package compton_pages

import (
	"github.com/arelate/gaugin/rest/compton_data"
	"github.com/arelate/gaugin/rest/compton_fragments"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/align"
	"github.com/boggydigital/compton/consts/color"
	"github.com/boggydigital/compton/consts/direction"
	"github.com/boggydigital/compton/consts/font_weight"
	"github.com/boggydigital/compton/consts/size"
	"github.com/boggydigital/compton/elements/details_summary"
	"github.com/boggydigital/compton/elements/els"
	"github.com/boggydigital/compton/elements/flex_items"
	"github.com/boggydigital/compton/elements/fspan"
	"github.com/boggydigital/compton/elements/svg_use"
	"github.com/boggydigital/kevlar"
)

var osOrder = []vangogh_local_data.OperatingSystem{
	vangogh_local_data.Windows,
	vangogh_local_data.MacOS,
	vangogh_local_data.Linux,
	vangogh_local_data.AnyOperatingSystem,
}

type DownloadVariant struct {
	dlType   vangogh_local_data.DownloadType
	version  string
	langCode string
}

var downloadTypesStrings = map[vangogh_local_data.DownloadType]string{
	vangogh_local_data.Installer: "Installer",
	vangogh_local_data.DLC:       "DLC",
	vangogh_local_data.Extra:     "Extra",
	vangogh_local_data.Movie:     "Movie",
}

// Downloads will present available installers, DLCs in the following hierarchy:
// - Operating system heading - Installers and DLCs (separately)
// - title_values list of downloads by version
func Downloads(id string, clientOs vangogh_local_data.OperatingSystem, dls vangogh_local_data.DownloadsList, rdx kevlar.ReadableRedux) compton.Element {
	s := compton_fragments.ProductSection(compton_data.DownloadsSection)

	pageStack := flex_items.FlexItems(s, direction.Column)
	s.Append(pageStack)

	dlOs := downloadsOperatingSystems(dls)

	for ii, os := range dlOs {
		osRow := flex_items.FlexItems(s, direction.Row).
			AlignItems(align.Center).
			ColumnGap(size.Small)
		osSymbol := svg_use.Sparkle
		if smb, ok := compton_data.OperatingSystemSymbols[os]; ok {
			osSymbol = smb
		}
		osIcon := svg_use.SvgUse(s, osSymbol)
		osTitle := els.H3()
		osString := ""
		switch os {
		case vangogh_local_data.AnyOperatingSystem:
			osString = "Goodies"
		default:
			osString = os.String()
		}
		osTitle.Append(fspan.Text(s, osString).ForegroundColor(color.Gray))
		osRow.Append(osIcon, osTitle)
		pageStack.Append(osRow)

		variants := getDownloadVariants(os, dls)
		for _, dv := range variants {

			row := flex_items.FlexItems(s, direction.Row).
				ColumnGap(size.Small).
				AlignItems(align.Center)

			typeIcon := svg_use.SvgUse(s, svg_use.Circle)
			typeIcon.AddClass(dv.dlType.String())
			row.Append(typeIcon)

			typeSpan := fspan.Text(s, downloadTypesStrings[dv.dlType]).FontWeight(font_weight.Bolder)
			row.Append(typeSpan)

			versionSpan := fspan.Text(s, dv.version).ForegroundColor(color.Gray)
			row.Append(versionSpan)

			lcSpan := fspan.Text(s, compton_data.LanguageFlags[dv.langCode])
			row.Append(lcSpan)

			dsDownloadVariant := details_summary.Smaller(s, row, os == clientOs)
			pageStack.Append(dsDownloadVariant)

			downloads := filterDownloads(os, dls, dv)

			downloadsRow := flex_items.FlexItems(s, direction.Row)
			dsDownloadVariant.Append(downloadsRow)

			for _, dl := range downloads {
				name := dl.Name
				if dl.Type == vangogh_local_data.DLC {
					name = dl.ProductTitle
				}
				link := els.A(dl.ManualUrl)
				linkTitle := fspan.Text(s, name).FontSize(size.Small).FontWeight(font_weight.Bolder)
				link.Append(linkTitle)
				link.AddClass("download")
				downloadsRow.Append(link)
			}
		}

		if ii != len(dlOs)-1 {
			pageStack.Append(els.Hr())
		}

	}

	return s
}

func downloadsOperatingSystems(dls vangogh_local_data.DownloadsList) []vangogh_local_data.OperatingSystem {
	dlOs := make(map[vangogh_local_data.OperatingSystem]any)
	for _, dl := range dls {
		dlOs[dl.OS] = nil
	}

	oses := make([]vangogh_local_data.OperatingSystem, 0, len(dlOs))
	for _, os := range osOrder {
		if _, ok := dlOs[os]; ok {
			oses = append(oses, os)
		}
	}
	return oses
}

func (dv *DownloadVariant) Equals(other *DownloadVariant) bool {
	return dv.dlType == other.dlType &&
		dv.version == other.version &&
		dv.langCode == other.langCode
}

func hasDownloadVariant(dvs []*DownloadVariant, other *DownloadVariant) bool {
	for _, dv := range dvs {
		if dv.Equals(other) {
			return true
		}
	}
	return false
}

func getDownloadVariants(os vangogh_local_data.OperatingSystem, dls vangogh_local_data.DownloadsList) []*DownloadVariant {

	variants := make([]*DownloadVariant, 0)
	for _, dl := range dls {
		if dl.OS != os {
			continue
		}

		dv := &DownloadVariant{
			dlType:   dl.Type,
			version:  dl.Version,
			langCode: dl.LanguageCode,
		}

		if !hasDownloadVariant(variants, dv) {
			variants = append(variants, dv)
		}

	}
	return variants
}

func filterDownloads(os vangogh_local_data.OperatingSystem, dls vangogh_local_data.DownloadsList, dv *DownloadVariant) []vangogh_local_data.Download {
	downloads := make([]vangogh_local_data.Download, 0)
	for _, dl := range dls {
		if dl.OS != os ||
			dl.Type != dv.dlType ||
			dv.version != dl.Version ||
			dv.langCode != dl.LanguageCode {
			continue
		}
		downloads = append(downloads, dl)
	}
	return downloads
}

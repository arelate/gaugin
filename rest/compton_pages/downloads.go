package compton_pages

import (
	"fmt"
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

	if _, ok := rdx.GetLastVal(vangogh_local_data.ValidationCompletedProperty, id); ok {
		if validationResults, sure := rdx.GetAllValues(vangogh_local_data.ValidationResultProperty, id); sure && len(validationResults) > 0 {

			lastResult := validationResults[len(validationResults)-1]
			validationSection := flex_items.FlexItems(s, direction.Row).JustifyContent(align.Center)
			validationSection.Append(els.DivText(lastResult))
			validationSection.AddClass("validation-results", lastResult)

			pageStack.Append(validationSection)
		}
	}

	dlOs := downloadsOperatingSystems(dls)

	for ii, os := range dlOs {

		if osHeading := operatingSystemHeading(s, os); osHeading != nil {
			pageStack.Append(osHeading)
		}

		variants := getDownloadVariants(os, dls)
		for _, variant := range variants {
			if dv := downloadVariant(s, variant); dv != nil {
				pageStack.Append(dv)
			}
			if dlLinks := downloadLinks(s, os, variant, dls); dlLinks != nil {
				pageStack.Append(dlLinks)

			}

		}

		if ii != len(dlOs)-1 {
			pageStack.Append(els.Hr())
		}
	}

	return s
}

func operatingSystemHeading(r compton.Registrar, os vangogh_local_data.OperatingSystem) compton.Element {
	osRow := flex_items.FlexItems(r, direction.Row).
		AlignItems(align.Center).
		ColumnGap(size.Small)
	osSymbol := svg_use.Sparkle
	if smb, ok := compton_data.OperatingSystemSymbols[os]; ok {
		osSymbol = smb
	}
	osIcon := svg_use.SvgUse(r, osSymbol)
	osIcon.AddClass("operating-system")
	osTitle := els.H3()
	osString := ""
	switch os {
	case vangogh_local_data.AnyOperatingSystem:
		osString = "Goodies"
	default:
		osString = os.String()
	}
	osTitle.Append(fspan.Text(r, osString))
	osRow.Append(osIcon, osTitle)
	return osRow
}

func downloadVariant(r compton.Registrar, dv *DownloadVariant) compton.Element {
	//column := flex_items.FlexItems(r, direction.Column).
	//	ColumnGap(size.XSmall).
	//	AlignItems(align.Start)

	row := flex_items.FlexItems(r, direction.Row).
		ColumnGap(size.Small).
		AlignItems(align.Center)

	typeIcon := svg_use.SvgUse(r, svg_use.Circle)
	typeIcon.AddClass(dv.dlType.String())
	typeSpan := fspan.Text(r, downloadTypesStrings[dv.dlType]).
		FontWeight(font_weight.Bolder).
		FontSize(size.Small)

	row.Append(typeSpan, typeIcon)

	if dv.langCode != "" {
		lcTitle := fspan.Text(r, "Lang:").FontSize(size.Small).ForegroundColor(color.Gray)
		lcSpan := fspan.Text(r, compton_data.LanguageFlags[dv.langCode])
		row.Append(lcTitle, lcSpan)
	}

	//column.Append(row)

	if dv.version != "" {
		versionTitle := fspan.Text(r, "Version:").FontSize(size.Small).ForegroundColor(color.Gray)
		versionSpan := fspan.Text(r, dv.version).FontSize(size.Small)
		row.Append(versionTitle, versionSpan)
	}

	return row
}

func downloadLinks(r compton.Registrar, os vangogh_local_data.OperatingSystem, dv *DownloadVariant, dls vangogh_local_data.DownloadsList) compton.Element {
	downloadLinksTitle := fspan.Text(r, "Show download links").FontWeight(font_weight.Bolder)

	dsDownloadLinks := details_summary.Smaller(r, downloadLinksTitle, false)

	downloads := filterDownloads(os, dls, dv)

	downloadsRow := flex_items.FlexItems(r, direction.Row).
		ColumnGap(size.Large).
		RowGap(size.Small)
	dsDownloadLinks.Append(downloadsRow)

	for _, dl := range downloads {
		if link := downloadLink(r, dl); link != nil {
			downloadsRow.Append(link)
		}
	}

	return dsDownloadLinks
}

func downloadLink(r compton.Registrar, dl vangogh_local_data.Download) compton.Element {

	link := els.A(dl.ManualUrl)
	link.AddClass("download")

	linkColumn := flex_items.FlexItems(r, direction.Column).
		RowGap(size.Unset)

	name := dl.Name
	if dl.Type == vangogh_local_data.DLC {
		name = dl.ProductTitle
	}
	linkTitle := fspan.Text(r, name).
		FontWeight(font_weight.Bolder)
	//ForegroundColor(color.Gray)
	linkColumn.Append(linkTitle)

	sizeRow := flex_items.FlexItems(r, direction.Row).
		ColumnGap(size.XSmall)
	sizeTitle := fspan.Text(r, "Size:").
		FontSize(size.Small).
		ForegroundColor(color.Gray)
	sizeSpan := fspan.Text(r, fmtBytes(dl.EstimatedBytes)).
		FontSize(size.Small)
	sizeRow.Append(sizeTitle, sizeSpan)
	linkColumn.Append(sizeRow)

	link.Append(linkColumn)

	return link
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

func fmtBytes(b int) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}

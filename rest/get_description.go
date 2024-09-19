package rest

import (
	"github.com/arelate/gaugin/paths"
	"github.com/arelate/gaugin/rest/compton_data"
	"github.com/boggydigital/compton/consts/color"
	"github.com/boggydigital/compton/consts/size"
	"github.com/boggydigital/compton/elements/els"
	"github.com/boggydigital/compton/elements/fspan"
	"github.com/boggydigital/compton/elements/iframe_expand"
	"github.com/boggydigital/compton/elements/recipes"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
)

func GetDescription(w http.ResponseWriter, r *http.Request) {

	// GET /description?id

	id := r.URL.Query().Get("id")

	idRedux, err := getRedux(
		http.DefaultClient,
		id,
		false,
		vangogh_local_data.DescriptionOverviewProperty,
		vangogh_local_data.DescriptionFeaturesProperty,
		vangogh_local_data.AdditionalRequirementsProperty,
		vangogh_local_data.CopyrightsProperty)

	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	rdx := idRedux[id]

	desc := ""
	if rdx != nil {
		if dop := rdx[vangogh_local_data.DescriptionOverviewProperty]; len(dop) > 0 {
			desc = dop[0]
		}
		if dfp := rdx[vangogh_local_data.DescriptionFeaturesProperty]; len(dfp) > 0 {
			desc += implicitToExplicitList(dfp[0])
		}
	}

	section := compton_data.DescriptionSection
	ifc := iframe_expand.IframeExpandContent(section, compton_data.SectionTitles[section])

	if desc == "" {
		fs := fspan.Text(ifc, "Description is not available for this product").
			ForegroundColor(color.Subtle)
		ifc.Append(recipes.Center(ifc, fs))
	} else {
		desc = rewriteItemsLinks(desc)
		desc = rewriteGameLinks(desc)
		desc = rewriteLinksAsTargetTop(desc)
		desc = fixQuotes(desc)
		desc = replaceDataFallbackUrls(desc)
		desc = rewriteVideoAsInline(desc)

		ifc.Append(els.Text(desc))
	}

	addtReqs := ""
	if rdx != nil {
		if arp := rdx[vangogh_local_data.AdditionalRequirementsProperty]; len(arp) > 0 {
			addtReqs = arp[0]
		}
	}
	if addtReqs != "" {
		arfs := fspan.Text(ifc, addtReqs).ForegroundColor(color.Subtle).FontSize(size.Small)
		ifc.Append(els.Br())
		ifc.Append(arfs)
	}

	copyright := ""
	if rdx != nil {
		if cp := rdx[vangogh_local_data.CopyrightsProperty]; len(cp) > 0 {
			copyright = cp[0]
		}
	}
	if copyright != "" {
		cfs := fspan.Text(ifc, copyright).ForegroundColor(color.Subtle).FontSize(size.Small)
		ifc.Append(els.Br())
		ifc.Append(cfs)
	}

	if err := ifc.WriteContent(w); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
	}
}

func rewriteItemsLinks(desc string) string {

	itemsUrls := vangogh_local_data.ExtractDescItems(desc)

	for _, itemUrl := range itemsUrls {
		if u, err := url.Parse(itemUrl); err != nil {
			continue
		} else {
			ggUrl := "/items" + u.Path
			desc = strings.Replace(desc, itemUrl, ggUrl, -1)
		}
	}

	return desc
}

func rewriteGameLinks(desc string) string {
	gameLinks := vangogh_local_data.ExtractGameLinks(desc)

	for _, gameLink := range gameLinks {
		if u, err := url.Parse(gameLink); err != nil {
			continue
		} else {
			_, slug := path.Split(u.Path)
			ggUrl := paths.ProductSlug(slug)
			desc = strings.Replace(desc, gameLink, ggUrl, -1)
		}
	}

	return desc
}

func rewriteLinksAsTargetTop(desc string) string {
	return strings.Replace(desc, "<a ", "<a target='_top' ", -1)
}

func rewriteVideoAsInline(desc string) string {
	return strings.Replace(desc, "<video ", "<video playsinline ", -1)
}

func fixQuotes(desc string) string {
	return strings.Replace(desc, "”", "\"", -1)
}

func replaceDataFallbackUrls(desc string) string {
	return strings.Replace(desc, "data-fallbackurl", "poster", -1)
}

const doubleNewLineChar = "\n\n"
const newLineChar = "\n"
const emDashCode = "\u2013"

// implicitToExplicitList looks for embedded characters
// that GOG.com is using for <ul> lists creation, e.g.
// https://www.gog.com/en/game/deaths_gambit
// and replaces that segment with explicit unordered lists.
// Currently known characters are listed as consts above this func.
func implicitToExplicitList(text string) string {
	var items []string
	if strings.Contains(text, doubleNewLineChar) {
		items = strings.Split(text, doubleNewLineChar)
	} else if strings.Contains(text, newLineChar) {
		items = strings.Split(text, newLineChar)
	} else if strings.Contains(text, emDashCode) {
		items = strings.Split(text, emDashCode)
	}

	if len(items) > 0 {
		builder := strings.Builder{}
		builder.WriteString("<ul>")
		for _, item := range items {
			builder.WriteString("<li>" + item + "</li>")
		}
		builder.WriteString("</ul>")
		text = builder.String()
	}

	return text
}

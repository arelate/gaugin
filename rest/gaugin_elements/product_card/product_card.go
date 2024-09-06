package product_card

import (
	"bytes"
	_ "embed"
	"github.com/arelate/gaugin/rest/compton_data"
	"github.com/arelate/gaugin/rest/gaugin_atoms"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/elements/els"
	"github.com/boggydigital/compton/elements/issa_image"
	"github.com/boggydigital/compton/elements/svg_inline"
	"github.com/boggydigital/issa"
	"github.com/boggydigital/kevlar"
	"io"
	"strings"
)

const (
	productCardElementName = "product-card"
)

var (
	//go:embed "markup/product-card.html"
	markupProductCard []byte
	//go:embed "style/product-card.css"
	ProductCardStyle []byte
)

var operatingSystemSymbols = map[vangogh_local_data.OperatingSystem]compton.Element{
	vangogh_local_data.Windows: svg_inline.SvgInline(svg_inline.Windows),
	vangogh_local_data.MacOS:   svg_inline.SvgInline(svg_inline.MacOS),
	vangogh_local_data.Linux:   svg_inline.SvgInline(svg_inline.Linux),
}

type ProductCardElement struct {
	compton.BaseElement
	wcr    compton.Registrar
	poster compton.Element
	rdx    kevlar.ReadableRedux
	id     string
}

func (pc *ProductCardElement) WriteRequirements(w io.Writer) error {
	//if pc.wcr.RequiresRegistration(productCardElementName) {
	//	if err := custom_elements.Define(w, custom_elements.Defaults(productCardElementName)); err != nil {
	//		return err
	//	}
	//	if _, err := io.Copy(w, bytes.NewReader(markupTemplate)); err != nil {
	//		return err
	//	}
	//}
	if pc.poster != nil {
		if err := pc.poster.WriteRequirements(w); err != nil {
			return err
		}
	}
	return pc.BaseElement.WriteRequirements(w)
}

func (pc *ProductCardElement) WriteContent(w io.Writer) error {
	return compton.WriteContents(bytes.NewReader(markupProductCard), w, pc.elementFragmentWriter)
}

func (pc *ProductCardElement) elementFragmentWriter(t string, w io.Writer) error {
	switch t {
	case ".Id":
		if _, err := io.WriteString(w, pc.id); err != nil {
			return err
		}
	case ".Poster":
		if pc.poster != nil {
			if err := pc.poster.WriteContent(w); err != nil {
				return err
			}
		}
	case ".Title":
		if title, ok := pc.rdx.GetLastVal(vangogh_local_data.TitleProperty, pc.id); ok {
			if _, err := io.WriteString(w, title); err != nil {
				return err
			}
		}
	case ".Labels":
		for _, property := range compton_data.LabelProperties {
			if title := compton_data.LabelTitle(pc.id, property, pc.rdx); title != "" {
				class := compton_data.LabelClass(pc.id, property, pc.rdx)
				label := createLabel(property, title, class)
				if err := label.WriteContent(w); err != nil {
					return err
				}
			}
		}
	case ".OperatingSystems":

		if oses, ok := pc.rdx.GetAllValues(vangogh_local_data.OperatingSystemsProperty, pc.id); ok {
			for _, os := range vangogh_local_data.ParseManyOperatingSystems(oses) {
				symbol := operatingSystemSymbols[os]
				if err := symbol.WriteContent(w); err != nil {
					return err
				}
			}
		}
	case ".Developers":
		if developers, ok := pc.rdx.GetAllValues(vangogh_local_data.DevelopersProperty, pc.id); ok {
			if _, err := io.WriteString(w, strings.Join(developers, ", ")); err != nil {
				return err
			}
		}
	case ".Publishers":
		if publishers, ok := pc.rdx.GetAllValues(vangogh_local_data.PublishersProperty, pc.id); ok {
			if _, err := io.WriteString(w, strings.Join(publishers, ", ")); err != nil {
				return err
			}
		}

	case compton.AttributesToken:
		return pc.BaseElement.WriteFragment(compton.AttributesToken, w)
	default:
		return compton.ErrUnknownToken(t)
	}
	return nil
}

func (pc *ProductCardElement) SetDehydratedPoster(dehydratedSrc, posterSrc string) *ProductCardElement {
	pc.poster = issa_image.IssaImageDehydrated(pc.wcr, dehydratedSrc, posterSrc)
	pc.poster.SetAttr("slot", "poster")
	return pc
}

func (pc *ProductCardElement) SetHydratedPoster(hydratedSrc, posterSrc string) *ProductCardElement {
	pc.poster = issa_image.IssaImageHydrated(pc.wcr, hydratedSrc, posterSrc)
	pc.poster.SetAttr("slot", "poster")
	return pc
}

func createLabel(property, title, class string) compton.Element {
	label := els.ListItemText(title)
	cs := []string{"label", property, title, class}
	label.SetClass(cs...)
	return label
}

func ProductCard(wcr compton.Registrar, id string, hydrated bool, rdx kevlar.ReadableRedux) *ProductCardElement {
	pc := &ProductCardElement{
		BaseElement: compton.BaseElement{
			TagName: gaugin_atoms.ProductCard,
			Markup:  markupProductCard,
		},
		wcr: wcr,
		id:  id,
		rdx: rdx,
	}

	if viSrc, ok := rdx.GetLastVal(vangogh_local_data.VerticalImageProperty, id); ok {
		dhSrc, _ := rdx.GetLastVal(vangogh_local_data.DehydratedVerticalImageProperty, id)
		if hydrated {
			hSrc := issa.HydrateColor(dhSrc)
			pc.SetHydratedPoster(hSrc, "/image?id="+viSrc)
		} else {
			pc.SetDehydratedPoster(dhSrc, "/image?id="+viSrc)
		}
	}

	pc.SetAttr("data-id", id)

	return pc
}

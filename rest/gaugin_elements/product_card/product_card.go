package product_card

import (
	"bytes"
	_ "embed"
	"github.com/arelate/gaugin/rest/gaugin_atoms"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/custom_elements"
	"github.com/boggydigital/compton/elements/issa_image"
	"github.com/boggydigital/compton/elements/svg_inline"
	"github.com/boggydigital/kevlar"
	"io"
	"strings"
)

const (
	productCardElementName = "product-card"
)

var (
	//go:embed "markup/template.html"
	markupTemplate []byte
	//go:embed "markup/product-card.html"
	markupProductCard []byte
)

var operatingSystemSymbols = map[vangogh_local_data.OperatingSystem]compton.Element{
	vangogh_local_data.Windows: svg_inline.New(svg_inline.Windows),
	vangogh_local_data.MacOS:   svg_inline.New(svg_inline.MacOS),
	vangogh_local_data.Linux:   svg_inline.New(svg_inline.Linux),
}

type ProductCard struct {
	compton.BaseElement
	wcr    compton.Registrar
	poster compton.Element
	//title            string
	//labels           []compton.Element
	//operatingSystems []compton.Element
	//developers       []string
	//publishers       []string
	rdx kevlar.ReadableRedux
	id  string
}

func (pc *ProductCard) WriteRequirements(w io.Writer) error {
	if pc.wcr.RequiresRegistration(productCardElementName) {
		if err := custom_elements.Define(w, custom_elements.Defaults(productCardElementName)); err != nil {
			return err
		}
		if _, err := io.Copy(w, bytes.NewReader(markupTemplate)); err != nil {
			return err
		}
	}
	if pc.poster != nil {
		if err := pc.poster.WriteRequirements(w); err != nil {
			return err
		}
	}
	return pc.BaseElement.WriteRequirements(w)
}

func (pc *ProductCard) WriteContent(w io.Writer) error {
	return compton.WriteContents(bytes.NewReader(markupProductCard), w, pc.elementFragmentWriter)
}

func (pc *ProductCard) elementFragmentWriter(t string, w io.Writer) error {
	switch t {
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
		//for _, label := range pc.labels {
		//	if err := label.WriteContent(w); err != nil {
		//		return err
		//	}
		//}
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

//	func (pc *ProductCard) SetDevelopers(developers ...string) *ProductCard {
//		pc.developers = developers
//		return pc
//	}
//
//	func (pc *ProductCard) SetPublishers(publishers ...string) *ProductCard {
//		pc.publishers = publishers
//		return pc
//	}
func (pc *ProductCard) SetDehydratedPoster(dehydratedSrc, posterSrc string) *ProductCard {
	pc.poster = issa_image.NewDehydrated(pc.wcr, dehydratedSrc, posterSrc)
	pc.poster.SetAttr("slot", "poster")
	return pc
}

func (pc *ProductCard) SetHydratedPoster(hydratedSrc, posterSrc string) *ProductCard {
	pc.poster = issa_image.NewHydrated(pc.wcr, hydratedSrc, posterSrc)
	pc.poster.SetAttr("slot", "poster")
	return pc
}

//func operatingSystemSymbol()

//func (pc *ProductCard) SetOperatingSystems(operatingSystems ...vangogh_local_data.OperatingSystem) *ProductCard {
//	pc.operatingSystems = nil
//	for _, os := range operatingSystems {
//		var symbol svg_inline.Symbol
//		switch os {
//		case vangogh_local_data.Windows:
//			symbol = svg_inline.Windows
//		case vangogh_local_data.MacOS:
//			symbol = svg_inline.MacOS
//		case vangogh_local_data.Linux:
//			symbol = svg_inline.Linux
//		default:
//			panic("unknown operating system")
//		}
//		pc.operatingSystems = append(pc.operatingSystems, svg_inline.New(symbol))
//	}
//	return pc
//}

//
//func (pc *ProductCard) SetLabels(values map[string]string, classes map[string][]string, order ...string) *ProductCard {
//
//	pc.labels = nil
//
//	if order == nil {
//		order = maps.Keys(values)
//		slices.Sort(order)
//	}
//
//	for _, l := range order {
//		label := els.NewDiv()
//
//		value := values[l]
//		label.Append(els.NewText(value))
//		cs := []string{"label", l, value}
//		if lcs, ok := classes[l]; ok {
//			cs = append(cs, lcs...)
//		}
//		label.SetClass(cs...)
//		pc.labels = append(pc.labels, label)
//	}
//
//	return pc
//}
//
//func (pc *ProductCard) SetTitle(title string) *ProductCard {
//	pc.title = title
//	return pc
//}

func New(wcr compton.Registrar, id string, rdx kevlar.ReadableRedux) *ProductCard {
	pc := &ProductCard{
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
		//hSrc := issa.HydrateColor(dhSrc)
		pc.SetDehydratedPoster(dhSrc, "image?id="+viSrc)
	}

	pc.SetAttr("data-id", id)

	return pc
}

//func NewData(wcr compton.Registrar, id string, rdx kevlar.ReadableRedux) (*ProductCard, error) {
//
//	if err := rdx.MustHave(
//		vangogh_local_data.DehydratedVerticalImageProperty,
//		vangogh_local_data.TitleProperty); err != nil {
//		return nil, err
//	}
//
//	pc := New(wcr, id)
//
//	if viSrc, ok := rdx.GetLastVal(vangogh_local_data.VerticalImageProperty, id); ok {
//		dhSrc, _ := rdx.GetLastVal(vangogh_local_data.DehydratedVerticalImageProperty, id)
//		pc.SetDehydratedPoster(dhSrc, "image?id="+viSrc)
//	}
//
//	if title, ok := rdx.GetLastVal(vangogh_local_data.TitleProperty, id); ok {
//		pc.SetTitle(title)
//	}
//
//	if oses, ok := rdx.GetAllValues(vangogh_local_data.OperatingSystemsProperty, id); ok {
//		poses := vangogh_local_data.ParseManyOperatingSystems(oses)
//		pc.SetOperatingSystems(poses...)
//	}
//
//	if developers, ok := rdx.GetAllValues(vangogh_local_data.DevelopersProperty, id); ok {
//		pc.SetDevelopers(developers...)
//	}
//
//	if publishers, ok := rdx.GetAllValues(vangogh_local_data.PublishersProperty, id); ok {
//		pc.SetPublishers(publishers...)
//	}
//
//	return pc, nil
//}

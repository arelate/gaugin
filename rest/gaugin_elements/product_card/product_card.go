package product_card

import (
	"bytes"
	_ "embed"
	"github.com/arelate/gaugin/rest/gaugin_atoms"
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/direction"
	"github.com/boggydigital/compton/consts/size"
	"github.com/boggydigital/compton/custom_elements"
	"github.com/boggydigital/compton/elements/els"
	"github.com/boggydigital/compton/elements/flex_items"
	"github.com/boggydigital/compton/elements/issa_image"
	"github.com/boggydigital/compton/elements/svg_inline"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
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

type ProductCard struct {
	compton.BaseElement
	wcr              compton.Registrar
	poster           compton.Element
	title            string
	labels           compton.Element
	operatingSystems compton.Element
	developers       []string
	publishers       []string
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
		if _, err := io.WriteString(w, pc.title); err != nil {
			return err
		}
	case ".Labels":
		if err := pc.labels.WriteContent(w); err != nil {
			return err
		}
	case ".OperatingSystems":
		if err := pc.operatingSystems.WriteContent(w); err != nil {
			return err
		}
	case ".Developers":
		if _, err := io.WriteString(w, strings.Join(pc.developers, ", ")); err != nil {
			return err
		}
	case ".Publishers":
		if _, err := io.WriteString(w, strings.Join(pc.publishers, ", ")); err != nil {
			return err
		}
	default:
		return compton.ErrUnknownToken(t)
	}
	return nil
}

func (pc *ProductCard) SetDevelopers(developers ...string) *ProductCard {
	pc.developers = developers
	return pc
}

func (pc *ProductCard) SetPublishers(publishers ...string) *ProductCard {
	pc.publishers = publishers
	return pc
}

func (pc *ProductCard) SetPoster(dehydratedSrc, posterSrc string) *ProductCard {
	pc.poster = issa_image.NewDehydrated(pc.wcr, dehydratedSrc, posterSrc)
	pc.poster.SetAttr("slot", "poster")
	return pc
}

func (pc *ProductCard) SetOperatingSystems(operatingSystems ...vangogh_local_data.OperatingSystem) *ProductCard {
	osFlexItems := flex_items.New(pc.wcr, direction.Row).
		SetColumnGap(size.Small)
	osFlexItems.SetAttr("slot", "operating-systems")
	for _, os := range operatingSystems {
		var symbol svg_inline.Symbol
		switch os {
		case vangogh_local_data.Windows:
			symbol = svg_inline.Windows
		case vangogh_local_data.MacOS:
			symbol = svg_inline.MacOS
		case vangogh_local_data.Linux:
			symbol = svg_inline.Linux
		default:
			panic("unknown operating system")
		}
		osFlexItems.Append(svg_inline.New(symbol))
	}
	pc.operatingSystems = osFlexItems
	return pc
}

func (pc *ProductCard) SetLabels(values map[string]string, classes map[string][]string, order ...string) *ProductCard {

	labelsFlexItems := flex_items.New(pc.wcr, direction.Row).
		SetColumnGap(size.XXSmall).
		SetRowGap(size.XXSmall)
	labelsFlexItems.SetClass("labels")
	labelsFlexItems.SetAttr("slot", "labels")

	if order == nil {
		order = maps.Keys(values)
		slices.Sort(order)
	}

	for _, l := range order {
		label := els.NewLabel("")
		value := values[l]
		label.Append(els.NewText(value))
		cs := []string{l, value}
		if lcs, ok := classes[l]; ok {
			cs = append(cs, lcs...)
		}
		label.SetClass(cs...)
		labelsFlexItems.Append(label)
	}

	pc.labels = labelsFlexItems
	return pc
}

func (pc *ProductCard) SetTitle(title string) *ProductCard {
	pc.title = title
	return pc
}

func New(wcr compton.Registrar) *ProductCard {
	return &ProductCard{
		BaseElement: compton.BaseElement{
			TagName: gaugin_atoms.ProductCard,
			Markup:  markupProductCard,
		},
		wcr: wcr,
	}
}

//func NewData(wcr compton.Registrar, id string, rdx kevlar.ReadableRedux) (*ProductCard, error) {
//
//	if err := rdx.MustHave(
//		vangogh_local_data.DehydratedVerticalImageProperty,
//		vangogh_local_data.TitleProperty); err != nil {
//		return nil, err
//	}
//
//	pc := New(wcr)
//
//	if viSrc, ok := rdx.GetLastVal(vangogh_local_data.VerticalImageProperty, id); ok {
//		dhSrc, _ := rdx.GetLastVal(vangogh_local_data.DehydratedVerticalImageProperty, id)
//		pc.SetPoster(dhSrc, viSrc)
//	}
//
//	if title, ok := rdx.GetLastVal(vangogh_local_data.TitleProperty, id); ok {
//		pc.SetTitle(title)
//	}
//
//	return pc, nil
//}

package product_card

import (
	"bytes"
	_ "embed"
	"github.com/arelate/gaugin/rest/gaugin_atoms"
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/custom_elements"
	"github.com/boggydigital/compton/elements/issa_image"
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
	wcr    compton.Registrar
	poster compton.Element
	title  string
	//labels           []string
	labels           compton.Element
	operatingSystems compton.Element
	//operatingSystems []vangogh_local_data.OperatingSystem
	developers []string
	publishers []string
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
	case ".OperatingSystems":
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

func New(wcr compton.Registrar, title string) *ProductCard {
	return &ProductCard{
		BaseElement: compton.BaseElement{
			TagName: gaugin_atoms.ProductCard,
			Markup:  markupProductCard,
		},
		wcr:   wcr,
		title: title,
	}
}

package product_labels

import (
	_ "embed"
	"fmt"
	"github.com/arelate/gaugin/rest/compton_data"
	"github.com/arelate/gaugin/rest/gaugin_atoms"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/class"
	"github.com/boggydigital/compton/consts/size"
	"github.com/boggydigital/compton/elements/els"
	"github.com/boggydigital/kevlar"
	"io"
)

const (
	registrationName      = "product-label"
	styleRegistrationName = "style-" + registrationName
)

var (
	//go:embed "style/product-labels.css"
	styleProductLabels []byte
)

type LabelsElement struct {
	compton.BaseElement
	r   compton.Registrar
	id  string
	rdx kevlar.ReadableRedux
	ul  compton.Element
}

func (lse *LabelsElement) WriteRequirements(w io.Writer) error {
	if lse.r.RequiresRegistration(styleRegistrationName) {
		if err := els.Style(styleProductLabels, styleRegistrationName).WriteContent(w); err != nil {
			return err
		}
	}
	if err := lse.ul.WriteRequirements(w); err != nil {
		return err
	}
	return lse.BaseElement.WriteRequirements(w)
}

func createLabel(title, property, class string) compton.Element {
	label := els.ListItemText(title)
	cs := []string{"label", property, title, class}
	label.AddClass(cs...)
	return label
}

func (lse *LabelsElement) WriteContent(w io.Writer) error {
	owned, _ := lse.rdx.GetLastVal(vangogh_local_data.OwnedProperty, lse.id)
	for _, property := range compton_data.LabelProperties {
		if vl, cl := labelValueClass(lse.id, property, owned, lse.rdx); vl != "" {
			lse.ul.Append(createLabel(vl, property, cl))
		}
	}
	return lse.ul.WriteContent(w)
}

func (lse *LabelsElement) FontSize(s size.Size) *LabelsElement {
	lse.ul.AddClass(class.FontSize(s))
	return lse
}

func Labels(r compton.Registrar, id string, rdx kevlar.ReadableRedux) *LabelsElement {
	return &LabelsElement{
		BaseElement: compton.BaseElement{
			TagName: gaugin_atoms.ProductLabels,
		},
		r:   r,
		id:  id,
		rdx: rdx,
		ul:  els.Ul(),
	}
}

func labelValueClass(id, property, owned string, rdx kevlar.ReadableRedux) (vl string, cl string) {
	vl, _ = rdx.GetLastVal(property, id)
	switch property {
	case vangogh_local_data.OwnedProperty:
		if res, ok := rdx.GetLastVal(vangogh_local_data.ValidationResultProperty, id); ok {
			if res == "OK" {
				cl = "validation-result-ok"
			} else {
				cl = "validation-result-err"
			}
		}
		fallthrough
	case vangogh_local_data.WishlistedProperty:
		fallthrough
	case vangogh_local_data.PreOrderProperty:
		fallthrough
	case vangogh_local_data.ComingSoonProperty:
		fallthrough
	case vangogh_local_data.InDevelopmentProperty:
		fallthrough
	case vangogh_local_data.IsFreeProperty:
		if vl == "true" {
			vl = compton_data.LabelTitles[property]
			break
		}
		vl = ""
	case vangogh_local_data.ProductTypeProperty:
		if vl == "GAME" {
			vl = ""
			break
		}
	case vangogh_local_data.DiscountPercentageProperty:
		if owned == "true" {
			vl = ""
			break
		}
		if vl != "" && vl != "0" {
			vl = fmt.Sprintf("-%s%%", vl)
			break
		}
		vl = ""
	case vangogh_local_data.TagIdProperty:
		if tagName, ok := rdx.GetLastVal(vangogh_local_data.TagNameProperty, vl); ok {
			vl = tagName
			break
		}
	case vangogh_local_data.DehydratedImageProperty:
		fallthrough
	case vangogh_local_data.DehydratedVerticalImageProperty:
		vl = property
	}
	return vl, cl
}

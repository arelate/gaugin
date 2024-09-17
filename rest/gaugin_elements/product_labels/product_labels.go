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
	for _, property := range compton_data.LabelProperties {
		if title := labelTitle(lse.id, property, lse.rdx); title != "" {
			class := labelClass(lse.id, property, lse.rdx)
			lse.ul.Append(createLabel(title, property, class))
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

func labelTitle(id, property string, rdx kevlar.ReadableRedux) string {

	label, _ := rdx.GetLastVal(property, id)
	owned, _ := rdx.GetLastVal(vangogh_local_data.OwnedProperty, id)

	switch property {
	case vangogh_local_data.WishlistedProperty:
		fallthrough
	case vangogh_local_data.OwnedProperty:
		fallthrough
	case vangogh_local_data.PreOrderProperty:
		fallthrough
	case vangogh_local_data.ComingSoonProperty:
		fallthrough
	case vangogh_local_data.InDevelopmentProperty:
		fallthrough
	case vangogh_local_data.IsFreeProperty:
		if label == "true" {
			return compton_data.LabelTitles[property]
		}
		return ""
	case vangogh_local_data.ProductTypeProperty:
		if label == "GAME" {
			return ""
		}
	case vangogh_local_data.DiscountPercentageProperty:
		if owned == "true" {
			return ""
		}
		if label != "" && label != "0" {
			return fmt.Sprintf("-%s%%", label)
		}
		return ""
	case vangogh_local_data.TagIdProperty:
		if tagName, ok := rdx.GetLastVal(vangogh_local_data.TagNameProperty, label); ok {
			return tagName
		}
	case vangogh_local_data.DehydratedImageProperty:
		fallthrough
	case vangogh_local_data.DehydratedVerticalImageProperty:
		return property
	}
	return label
}

func labelClass(id, property string, rdx kevlar.ReadableRedux) string {
	switch property {
	case vangogh_local_data.OwnedProperty:
		if res, ok := rdx.GetLastVal(vangogh_local_data.ValidationResultProperty, id); ok {
			if res == "OK" {
				return "validation-result-ok"
			} else {
				return "validation-result-err"
			}
		} else {
			return ""
		}
	}
	return ""
}

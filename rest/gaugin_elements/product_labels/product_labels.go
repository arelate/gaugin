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
	owned := false
	if op, ok := lse.rdx.GetLastVal(vangogh_local_data.OwnedProperty, lse.id); ok {
		owned = op == vangogh_local_data.TrueValue
	}
	for _, property := range compton_data.LabelProperties {
		if fmtLabel := formatLabel(lse.id, property, owned, lse.rdx); fmtLabel.value != "" {
			lse.ul.Append(createLabel(fmtLabel.value, property, fmtLabel.class))
		}
	}
	return lse.ul.WriteContent(w)
}

func (lse *LabelsElement) FontSize(s size.Size) *LabelsElement {
	lse.ul.AddClass(class.FontSize(s))
	return lse
}

func (lse *LabelsElement) RowGap(s size.Size) *LabelsElement {
	lse.ul.AddClass(class.RowGap(s))
	return lse
}

func (lse *LabelsElement) ColumnGap(s size.Size) *LabelsElement {
	lse.ul.AddClass(class.ColumnGap(s))
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

type formattedLabel struct {
	value string
	class string
}

func formatLabel(id, property string, owned bool, rdx kevlar.ReadableRedux) formattedLabel {

	fmtLabel := formattedLabel{}

	fmtLabel.value, _ = rdx.GetLastVal(property, id)
	switch property {
	case vangogh_local_data.OwnedProperty:
		if res, ok := rdx.GetLastVal(vangogh_local_data.ValidationResultProperty, id); ok {
			if res == "OK" {
				fmtLabel.class = "validation-result-ok"
			} else {
				fmtLabel.class = "validation-result-err"
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
		if fmtLabel.value == "true" {
			fmtLabel.value = compton_data.LabelTitles[property]
			break
		}
		fmtLabel.value = ""
	case vangogh_local_data.ProductTypeProperty:
		if fmtLabel.value == "GAME" {
			fmtLabel.value = ""
			break
		}
	case vangogh_local_data.DiscountPercentageProperty:
		if owned {
			fmtLabel.value = ""
			break
		}
		if fmtLabel.value != "" && fmtLabel.value != "0" {
			fmtLabel.value = fmt.Sprintf("-%s%%", fmtLabel.value)
			break
		}
		fmtLabel.value = ""
	case vangogh_local_data.TagIdProperty:
		if tagName, ok := rdx.GetLastVal(vangogh_local_data.TagNameProperty, fmtLabel.value); ok {
			fmtLabel.value = tagName
			break
		}
	case vangogh_local_data.DehydratedImageProperty:
		fallthrough
	case vangogh_local_data.DehydratedVerticalImageProperty:
		fmtLabel.value = property
	}
	return fmtLabel
}

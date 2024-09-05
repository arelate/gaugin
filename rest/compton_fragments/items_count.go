package compton_fragments

import (
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/alignment"
	"github.com/boggydigital/compton/elements/els"
	"github.com/boggydigital/compton/elements/flex_items"
	"strconv"
	"strings"
)

const (
	singleItem          = "1 item"
	manyItemsSinglePage = "{total} items"
	manyItemsManyPages  = "{from}-{to} out of {total} items"
)

func ItemsCount(r compton.Registrar, from, to, total int) compton.Element {
	title := ""
	switch total {
	case 1:
		title = singleItem
	case to:
		title = strings.Replace(manyItemsSinglePage, "{total}", strconv.Itoa(total), 1)
	default:
		title = strings.Replace(manyItemsManyPages, "{from}", strconv.Itoa(from+1), 1)
		title = strings.Replace(title, "{to}", strconv.Itoa(to), 1)
		title = strings.Replace(title, "{total}", strconv.Itoa(total), 1)
	}

	row := flex_items.FlexItemsRow(r).JustifyContent(alignment.Center)

	span := els.SpanText(title)
	span.SetClass("fg-subtle", "fs-xs")
	row.Append(span)

	return row
}

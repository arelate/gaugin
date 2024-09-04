package compton_pages

import (
	"github.com/arelate/gaugin/paths"
	"github.com/arelate/gaugin/rest/compton_data"
	"github.com/arelate/gaugin/rest/compton_fragments"
	"github.com/arelate/gaugin/rest/gaugin_elements/product_card"
	"github.com/arelate/gaugin/rest/gaugin_styles"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/direction"
	"github.com/boggydigital/compton/elements/details_toggle"
	"github.com/boggydigital/compton/elements/els"
	"github.com/boggydigital/compton/elements/flex_items"
	"github.com/boggydigital/compton/elements/grid_items"
	"github.com/boggydigital/compton/elements/page"
)

func SearchNew(query map[string][]string) compton.Element {

	p := page.New("Search - gaugin").
		SetFavIconEmoji("🪸").
		SetCustomStyles(gaugin_styles.GauginStyle)

	pageStack := flex_items.New(p, direction.Column)
	p.Append(pageStack)

	appNavLinks := compton_fragments.AppNavLinks(p, compton_data.AppNavSearch)
	searchLinks := compton_fragments.SearchLinks(p, compton_data.SearchNew)

	filterSearchDetails := details_toggle.NewClosed(p, "Filter & Search")

	filterSearchDetails.Append(compton_fragments.SearchForm(p))

	productCards := grid_items.New(p)

	testProductCard := product_card.New(p, "1965670180").
		SetTitle("SteamWorld Heist II").
		SetDevelopers("Thunderful Development").
		SetPublishers("Thunderful Publishing").
		SetPoster(
			"21=30=I/wATKIBAA4LABAgRKljIsOFCCA4f0qAxo2BDCBgnQsxI0cqTJytmYHxYrly8LhVplOtCDpyVcOFWgJtYkIY7k126xNvJE92MclaqALIiEgLOk+TIsQTHrsqfKqP+jApZsBy5HOCqVMmDZ9QoPFVSgcvD7g+CigTLzSiL76uoPKPy5JGRpwoeGSxo6tDBTtQoUXjA1tX6p3AMB3v3ehmntUqMxw0OHHDQIMbkxOVkgJNxYMUBGUnLefGSyAuCcjt26JjxxIoMFSo+eyHHBV6Z1KmDqNaxQoVOluXK5Fix5Z2OIMjRtSODnAQJplbKSEfnBE896WXIVB/F3MkKK03JkP/JDs6JnHjw4qFb8cQJoDJPNstsByhR+zxI0lUhV25POwNdlKFUPFrVE8gcSOCRhxxOtGNFPoGsIEM7h+jUTjrlAUKCHJbIgQQS7egzSj766MNeK+04kc4TJHyEBIMktMNKPeyQUgo7T6gQSBdPhMOCc4GA85EVrZTCST2AsMPeE4k8EmRM4OjDDlmttJJHJ1Y4aEWWiSSCSDtXXOFOIPqlE4ge5qWjU5dsIuLmTm52saIMM2hmBSJsikcGInrG4xELZTyyWnngiMenm26iQ0Z5TliBDjoercDCCl2skQgZXtbTSSf15KGCE53k0V45QTxQJxCPYPoIIvFMmccKZ8lDAFKgZJiawyO4OmIoGTk9YAFvK1QgnQ4P5KCom2QAAQQZiqLjwAMIBPDsAw44UEEO5WD7qKJAoANEDtBOVgG1DigaEAA7",
			"/image?id=3e352d9097ffd33cff8c43c7ddd31c2bb3d3a8f11089a6b737e21235777f71ad").
		SetOperatingSystems(vangogh_local_data.Windows).
		SetLabels(
			map[string]string{
				"owned": "Own",
			},
			map[string][]string{
				"owned": {"validation-result-ok"},
			})

	productLink := els.NewA(paths.ProductId("1965670180"))
	productLink.Append(testProductCard)

	for range 100 {

		productCards.Append(productLink)
	}

	pageStack.Append(
		appNavLinks,
		searchLinks,
		filterSearchDetails,
		productCards)

	return p
}

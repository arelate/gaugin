package rest

import (
	"crypto/sha256"
	"encoding/gob"
	"github.com/arelate/gog_integration"
	"github.com/arelate/steam_integration"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/middleware"
	"html/template"
	"io/fs"
)

var (
	tmpl             *template.Template
	operatingSystems []vangogh_local_data.OperatingSystem
	languageCodes    []string
)

func SetDownloadsOperatingSystems(os []vangogh_local_data.OperatingSystem) {
	operatingSystems = os
}

func SetDownloadsLanguageCodes(lc []string) {
	languageCodes = lc
}

func SetUsername(u string) {
	middleware.SetUsername(sha256.Sum256([]byte(u)))
}

func SetPassword(p string) {
	middleware.SetPassword(sha256.Sum256([]byte(p)))
}

func Init(templatesFS fs.FS) error {

	//GOG.com types
	gob.Register(gog_integration.AccountPage{})
	gob.Register(gog_integration.AccountProduct{})
	gob.Register(gog_integration.ApiProductV1{})
	gob.Register(gog_integration.ApiProductV2{})
	gob.Register(gog_integration.Details{})
	gob.Register(gog_integration.Licences{})
	gob.Register(gog_integration.OrderPage{})
	gob.Register(gog_integration.Order{})
	gob.Register(gog_integration.StorePage{})
	gob.Register(gog_integration.StoreProduct{})
	gob.Register(gog_integration.WishlistPage{})
	//Steam types
	gob.Register(steam_integration.AppList{})
	gob.Register(steam_integration.GetNewsForAppResponse{})
	gob.Register(steam_integration.AppReviews{})

	tmpl = template.Must(
		template.
			New("").
			Funcs(funcMap()).
			ParseFS(templatesFS, "templates/*.gohtml"))

	return nil
}

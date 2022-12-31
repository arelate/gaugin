package rest

import (
	"crypto/sha256"
	"encoding/gob"
	"github.com/arelate/gaugin/stencil_app"
	"github.com/arelate/gaugin/view_models"
	"github.com/arelate/southern_light/gog_integration"
	"github.com/arelate/southern_light/steam_integration"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/middleware"
	"github.com/boggydigital/stencil"
	"html/template"
	"io/fs"
)

const (
	AdminRole  = "admin"
	SharedRole = "shared"
)

var (
	tmpl             *template.Template
	operatingSystems []vangogh_local_data.OperatingSystem
	languageCodes    []string
	app              *stencil.AppConfiguration
)

func SetDownloadsOperatingSystems(os []vangogh_local_data.OperatingSystem) {
	operatingSystems = os
}

func SetDownloadsLanguageCodes(lc []string) {
	languageCodes = lc
}

func SetUsername(role, u string) {
	middleware.SetUsername(role, sha256.Sum256([]byte(u)))
}

func SetPassword(role, p string) {
	middleware.SetPassword(role, sha256.Sum256([]byte(p)))
}

func Init(templatesFS fs.FS, stencilAppStyles fs.FS) error {

	//GOG.com types
	gob.Register(gog_integration.AccountPage{})
	gob.Register(gog_integration.AccountProduct{})
	gob.Register(gog_integration.ApiProductV1{})
	gob.Register(gog_integration.ApiProductV2{})
	gob.Register(gog_integration.CatalogPage{})
	gob.Register(gog_integration.CatalogProduct{})
	gob.Register(gog_integration.Details{})
	gob.Register(gog_integration.Licences{})
	gob.Register(gog_integration.OrderPage{})
	gob.Register(gog_integration.Order{})
	gob.Register(gog_integration.UserWishlist{})
	//Steam types
	gob.Register(steam_integration.AppList{})
	gob.Register(steam_integration.GetNewsForAppResponse{})
	gob.Register(steam_integration.AppReviews{})

	tmpl = template.Must(
		template.
			New("").
			Funcs(view_models.FuncMap()).
			ParseFS(templatesFS, "templates/*.gohtml"))

	stencil.InitAppTemplates(stencilAppStyles, "stencil_app/styles/css.gohtml")

	var err error
	app, err = stencil_app.Init()

	return err
}

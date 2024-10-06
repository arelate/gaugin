package rest

import (
	"crypto/sha256"
	"encoding/gob"
	"github.com/arelate/southern_light/gog_integration"
	"github.com/arelate/southern_light/steam_integration"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/middleware"
)

const (
	AdminRole  = "admin"
	SharedRole = "shared"

	SearchResultsLimit = 60 // divisible by 2,3,4,5,6
)

var (
	operatingSystems []vangogh_local_data.OperatingSystem
	languageCodes    []string
	excludePatches   bool
)

func SetDefaultDownloadsFilters(
	os []vangogh_local_data.OperatingSystem,
	lc []string,
	ep bool) {
	operatingSystems = os
	languageCodes = lc
	excludePatches = ep
}

func SetUsername(role, u string) {
	middleware.SetUsername(role, sha256.Sum256([]byte(u)))
}

func SetPassword(role, p string) {
	middleware.SetPassword(role, sha256.Sum256([]byte(p)))
}

func Init() {

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
	gob.Register(steam_integration.DeckAppCompatibilityReport{})
	gob.Register(steam_integration.AppReviews{})
}

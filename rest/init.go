package rest

import (
	"crypto/sha256"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/kvas"
	"html/template"
	"io/fs"
)

var (
	tmpl             *template.Template
	operatingSystems []vangogh_local_data.OperatingSystem
	languageCodes    []string
	rxa              kvas.ReduxAssets
	usernameHash     [32]byte
	passwordHash     [32]byte
)

func SetDownloadsOperatingSystems(os []vangogh_local_data.OperatingSystem) {
	operatingSystems = os
}

func SetDownloadsLanguageCodes(lc []string) {
	languageCodes = lc
}

func SetUsername(u string) {
	usernameHash = sha256.Sum256([]byte(u))
}

func SetPassword(p string) {
	passwordHash = sha256.Sum256([]byte(p))
}

func Init(templatesFS fs.FS) error {
	var err error
	if rxa, err = vangogh_local_data.ConnectReduxAssets(
		vangogh_local_data.LocalManualUrlProperty); err != nil {
		return err
	}

	tmpl = template.Must(
		template.
			New("").
			Funcs(funcMap()).
			ParseFS(templatesFS, "templates/*.gohtml"))

	return nil
}

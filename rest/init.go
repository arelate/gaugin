package rest

import (
	"crypto/sha256"
	"github.com/arelate/vangogh_local_data"
	"html/template"
	"io/fs"
)

var (
	tmpl             *template.Template
	operatingSystems []vangogh_local_data.OperatingSystem
	languageCodes    []string
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
	tmpl = template.Must(
		template.
			New("").
			Funcs(funcMap()).
			ParseFS(templatesFS, "templates/*.gohtml"))

	return nil
}

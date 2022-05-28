package rest

import (
	"crypto/sha256"
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
	tmpl = template.Must(
		template.
			New("").
			Funcs(funcMap()).
			ParseFS(templatesFS, "templates/*.gohtml"))

	return nil
}

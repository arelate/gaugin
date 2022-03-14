package api

import (
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
)

func SetDownloadsOperatingSystems(os []vangogh_local_data.OperatingSystem) {
	operatingSystems = os
}

func SetDownloadsLanguageCodes(lc []string) {
	languageCodes = lc
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

package api

import (
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/kvas"
	"html/template"
	"io/fs"
)

var (
	tmpl             *template.Template
	cssFiles         fs.FS
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

func Init(htmlFS fs.FS, cssFS fs.FS) error {
	cssFiles = cssFS

	var err error
	if rxa, err = vangogh_local_data.ConnectReduxAssets(
		vangogh_local_data.LocalManualUrlProperty); err != nil {
		return err
	}

	tmpl = template.Must(
		template.
			New("").
			Funcs(funcMap()).
			ParseFS(htmlFS, "html/*.gohtml"))

	return nil
}

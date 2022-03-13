package api

import (
	"html/template"
	"io/fs"
)

var (
	tmpl     *template.Template
	cssFiles fs.FS
)

func Init(htmlFS fs.FS, cssFS fs.FS) {
	cssFiles = cssFS

	tmpl = template.Must(
		template.
			New("").
			Funcs(funcMap()).
			ParseFS(htmlFS, "html/*.gohtml"))
}

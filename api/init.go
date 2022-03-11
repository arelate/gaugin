package api

import (
	"html/template"
	"io/fs"
)

var tmpl *template.Template

func Init(templatesFS fs.FS) {
	tmpl = template.Must(template.New("").Funcs(funcMap()).ParseFS(templatesFS, "html/*.gohtml"))
}

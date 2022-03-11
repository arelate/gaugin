package main

import (
	"embed"
	"fmt"
	"github.com/arelate/gaugin/api"
	"log"
	"net/http"
	"sync"
)

var once = sync.Once{}

//go:embed "html/*.gohtml"
var htmlTemplates embed.FS

func main() {

	once.Do(func() { api.Init(htmlTemplates) })
	api.HandleFuncs()

	port := 1848
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatalln(err)
	}
}

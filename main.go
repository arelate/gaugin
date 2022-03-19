package main

import (
	"bytes"
	"embed"
	"github.com/arelate/gaugin/cli"
	"github.com/arelate/gaugin/rest"
	"github.com/boggydigital/clo"
	"log"
	"os"
	"sync"
)

var (
	once = sync.Once{}
	//go:embed "templates/*.gohtml"
	templates embed.FS
	//go:embed "cli-commands.txt"
	cliCommands []byte
	//go:embed "cli-help.txt"
	cliHelp []byte
)

func main() {

	once.Do(func() {
		if err := rest.Init(templates); err != nil {
			log.Fatalln(err)
		}
	})

	defs, err := clo.Load(
		bytes.NewBuffer(cliCommands),
		bytes.NewBuffer(cliHelp),
		nil)
	if err != nil {
		log.Fatalln(err)
	}

	clo.HandleFuncs(map[string]clo.Handler{
		"serve":   cli.ServeHandler,
		"version": cli.VersionHandler,
	})

	if err := defs.AssertCommandsHaveHandlers(); err != nil {
		log.Fatalln(err)
	}

	if err := defs.Serve(os.Args[1:]); err != nil {
		log.Fatalln(err)
	}
}

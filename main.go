package main

import (
	"bytes"
	"embed"
	"github.com/arelate/gaugin/cli"
	"github.com/arelate/gaugin/rest"
	"github.com/boggydigital/clo"
	"github.com/boggydigital/nod"
	"os"
	"sync"
)

var (
	once = sync.Once{}
	//go:embed "templates/*.gohtml"
	templates embed.FS
	//go:embed "stencil_app/styles/css.gohtml"
	stencilAppStyles embed.FS
	//go:embed "cli-commands.txt"
	cliCommands []byte
	//go:embed "cli-help.txt"
	cliHelp []byte
)

func main() {

	nod.EnableStdOutPresenter()

	gs := nod.Begin("gaugin is serving vangogh data")
	defer gs.End()

	once.Do(func() {
		if err := rest.Init(templates, stencilAppStyles); err != nil {
			_ = gs.EndWithError(err)
			os.Exit(1)
		}
	})

	defs, err := clo.Load(
		bytes.NewBuffer(cliCommands),
		bytes.NewBuffer(cliHelp),
		nil)
	if err != nil {
		_ = gs.EndWithError(err)
		os.Exit(1)
	}

	clo.HandleFuncs(map[string]clo.Handler{
		"serve":   cli.ServeHandler,
		"version": cli.VersionHandler,
	})

	if err := defs.AssertCommandsHaveHandlers(); err != nil {
		_ = gs.EndWithError(err)
		os.Exit(1)
	}

	if err := defs.Serve(os.Args[1:]); err != nil {
		_ = gs.EndWithError(err)
		os.Exit(1)
	}
}

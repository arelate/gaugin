package main

import (
	"bytes"
	_ "embed"
	"github.com/arelate/gaugin/cli"
	"github.com/arelate/gaugin/rest"
	"github.com/boggydigital/clo"
	"github.com/boggydigital/nod"
	"os"
	"sync"
)

var (
	once = sync.Once{}
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
		rest.Init()
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

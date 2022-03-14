package main

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/arelate/gaugin/api"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/clo"
	"github.com/boggydigital/nod"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
)

const defaultVangoghStateDir = "/var/lib/vangogh"

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

	defs, err := clo.Load(
		bytes.NewBuffer(cliCommands),
		bytes.NewBuffer(cliHelp),
		nil)
	if err != nil {
		log.Fatalln(err)
	}

	clo.HandleFuncs(map[string]clo.Handler{
		"serve": ServeHandler,
	})

	if err := defs.AssertCommandsHaveHandlers(); err != nil {
		log.Fatalln(err)
	}

	if err := defs.Serve(os.Args[1:]); err != nil {
		log.Fatalln(err)
	}
}

func ServeHandler(u *url.URL) error {
	portStr := vangogh_local_data.ValueFromUrl(u, "port")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return err
	}

	vangoghScheme := vangogh_local_data.ValueFromUrl(u, "vangogh_scheme")
	vangoghAddress := vangogh_local_data.ValueFromUrl(u, "vangogh_address")
	vangoghPortStr := vangogh_local_data.ValueFromUrl(u, "vangogh_port")
	vangoghPort, err := strconv.Atoi(vangoghPortStr)
	if err != nil {
		return err
	}

	if vangoghScheme == "" {
		return fmt.Errorf("missing vangogh scheme")
	}
	if vangoghAddress == "" {
		return fmt.Errorf("missing vangogh address")
	}
	if vangoghPortStr == "" {
		return fmt.Errorf("missing vangogh port")
	}

	api.SetVangoghConnection(vangoghScheme, vangoghAddress, vangoghPort)

	vangoghStateDir := vangogh_local_data.ValueFromUrl(u, "vangogh_state_directory")
	if vangoghStateDir == "" {
		vangoghStateDir = defaultVangoghStateDir
	}

	vangogh_local_data.ChRoot(vangoghStateDir)

	osStrings := vangogh_local_data.ValuesFromUrl(u, "operating_system")
	os := vangogh_local_data.ParseManyOperatingSystems(osStrings)
	lc := vangogh_local_data.ValuesFromUrl(u, "language_code")

	if len(os) == 0 {
		os = []vangogh_local_data.OperatingSystem{vangogh_local_data.AnyOperatingSystem}
	}
	if len(lc) == 0 {
		lc = []string{"en"}
	}

	api.SetDownloadsOperatingSystems(os)
	api.SetDownloadsLanguageCodes(lc)

	username := vangogh_local_data.ValueFromUrl(u, "username")
	password := vangogh_local_data.ValueFromUrl(u, "password")

	api.SetUsername(username)
	api.SetPassword(password)

	return Serve(port, vangogh_local_data.FlagFromUrl(u, "stderr"))
}

func Serve(port int, stderr bool) error {

	if stderr {
		nod.EnableStdErrLogger()
		nod.DisableOutput(nod.StdOut)
	}

	once.Do(func() {
		if err := api.Init(templates); err != nil {
			log.Fatalln(err)
		}
	})
	api.HandleFuncs()

	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

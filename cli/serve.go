package cli

import (
	"fmt"
	"github.com/arelate/gaugin/rest"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/http"
	"net/url"
	"strconv"
)

const (
	defaultRootDir          = "/var/lib/vangogh"
	defaultVangoghVideosDir = defaultRootDir + "/videos"
	defaultVangoghImagesDir = defaultRootDir + "/images"
	defaultVangoghItemsDir  = defaultRootDir + "/items"
)

func ServeHandler(u *url.URL) error {
	portStr := vangogh_local_data.ValueFromUrl(u, "port")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return err
	}

	vangoghScheme := vangogh_local_data.ValueFromUrl(u, "vangogh-scheme")
	vangoghAddress := vangogh_local_data.ValueFromUrl(u, "vangogh-address")
	vangoghPortStr := vangogh_local_data.ValueFromUrl(u, "vangogh-port")
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

	rest.SetVangoghConnection(vangoghScheme, vangoghAddress, vangoghPort)

	vangoghImagesDir := vangogh_local_data.ValueFromUrl(u, "vangogh-images-dir")
	if vangoghImagesDir == "" {
		vangoghImagesDir = defaultVangoghImagesDir
	}
	vangogh_local_data.SetImagesDir(vangoghImagesDir)

	vangoghItemsDir := vangogh_local_data.ValueFromUrl(u, "vangogh-items-dir")
	if vangoghItemsDir == "" {
		vangoghItemsDir = defaultVangoghItemsDir
	}
	vangogh_local_data.SetItemsDir(vangoghItemsDir)

	vangoghVideosDir := vangogh_local_data.ValueFromUrl(u, "vangogh-videos-dir")
	if vangoghVideosDir == "" {
		vangoghVideosDir = defaultVangoghVideosDir
	}
	vangogh_local_data.SetVideosDir(vangoghVideosDir)

	osStrings := vangogh_local_data.ValuesFromUrl(u, "operating-system")
	os := vangogh_local_data.ParseManyOperatingSystems(osStrings)
	lc := vangogh_local_data.ValuesFromUrl(u, "language-code")

	if len(os) == 0 {
		os = []vangogh_local_data.OperatingSystem{vangogh_local_data.AnyOperatingSystem}
	}
	if len(lc) == 0 {
		lc = []string{"en"}
	}

	rest.SetDownloadsOperatingSystems(os)
	rest.SetDownloadsLanguageCodes(lc)

	sharedUsername := vangogh_local_data.ValueFromUrl(u, "shared-username")
	sharedPassword := vangogh_local_data.ValueFromUrl(u, "shared-password")
	adminUsername := vangogh_local_data.ValueFromUrl(u, "admin-username")
	adminPassword := vangogh_local_data.ValueFromUrl(u, "admin-password")

	rest.SetUsername(rest.SharedRole, sharedUsername)
	rest.SetPassword(rest.SharedRole, sharedPassword)
	rest.SetUsername(rest.AdminRole, adminUsername)
	rest.SetPassword(rest.AdminRole, adminPassword)

	return Serve(port, vangogh_local_data.FlagFromUrl(u, "stderr"))
}

// Serve starts a web server, listening to the specified port with optional logging
func Serve(port int, stderr bool) error {

	if stderr {
		nod.EnableStdErrLogger()
		nod.DisableOutput(nod.StdOut)
	}

	rest.HandleFuncs(port)

	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

package cli

import (
	"fmt"
	"net/url"
	"runtime/debug"
)

func VersionHandler(u *url.URL) error {
	Version()
	return nil
}

func Version() {
	if bi, ok := debug.ReadBuildInfo(); ok {
		fmt.Println(bi)
	} else {
		fmt.Println("unable to read build info")
	}
}

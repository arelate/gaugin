package cli

import (
	"fmt"
	"net/url"
)

var (
	GitTag string
)

func VersionHandler(u *url.URL) error {
	Version()
	return nil
}

func Version() {
	if GitTag == "" {
		fmt.Println("version unknown")
	} else {
		fmt.Println(GitTag)
	}
}

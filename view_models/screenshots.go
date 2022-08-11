package view_models

import "github.com/arelate/vangogh_local_data"

type screenshots struct {
	Context     string
	Screenshots []string
}

func NewScreenshots(rdx map[string][]string) *screenshots {
	return &screenshots{
		Context:     "iframe",
		Screenshots: propertiesFromRedux(rdx, vangogh_local_data.ScreenshotsProperty),
	}
}

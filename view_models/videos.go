package view_models

import (
	"github.com/arelate/vangogh_local_data"
	"strings"
)

type videos struct {
	Context      string
	LocalVideos  []string
	RemoteVideos []string
}

func NewVideos(rdx map[string][]string) *videos {
	vvm := &videos{
		Context:      "iframe",
		LocalVideos:  make([]string, 0),
		RemoteVideos: make([]string, 0),
	}

	// filter videos to distinguish between locally available and remote videos

	for _, v := range propertiesFromRedux(rdx, vangogh_local_data.VideoIdProperty) {
		if !strings.Contains(v, "(") {
			vvm.LocalVideos = append(vvm.LocalVideos, v)
		} else {
			vvm.RemoteVideos = append(vvm.RemoteVideos, v)
		}
	}

	return vvm
}

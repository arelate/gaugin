package view_models

import (
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/kevlar"
)

type videos struct {
	Context      string
	LocalVideos  []string
	RemoteVideos []string
}

func NewVideos(id string, rdx kevlar.ReadableRedux) *videos {
	vvm := &videos{
		Context:      "iframe",
		LocalVideos:  make([]string, 0),
		RemoteVideos: make([]string, 0),
	}

	// filter videos to distinguish between locally available and remote videos

	videoIds, ok := rdx.GetAllValues(vangogh_local_data.VideoIdProperty, id)
	if !ok {
		return vvm
	}

	for _, v := range videoIds {

		if rdx.HasKey(vangogh_local_data.MissingVideoUrlProperty, v) {
			vvm.RemoteVideos = append(vvm.RemoteVideos, v)
		} else {
			vvm.LocalVideos = append(vvm.LocalVideos, v)
		}
	}

	return vvm
}

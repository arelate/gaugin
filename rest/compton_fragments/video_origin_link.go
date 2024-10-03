package compton_fragments

import (
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/color"
	"github.com/boggydigital/compton/consts/font_weight"
	"github.com/boggydigital/compton/consts/size"
	"github.com/boggydigital/compton/elements/els"
	"github.com/boggydigital/compton/elements/fspan"
	"github.com/boggydigital/yet_urls/youtube_urls"
)

func VideoOriginLink(r compton.Registrar, videoId string) compton.Element {
	originLink := els.A(youtube_urls.VideoUrl(videoId).String())
	originLink.SetAttribute("target", "_top")
	linkText := fspan.Text(r, "Watch at origin").
		FontSize(size.Small).
		FontWeight(font_weight.Bolder).
		ForegroundColor(color.Cyan)
	originLink.Append(linkText)

	return originLink
}

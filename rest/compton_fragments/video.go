package compton_fragments

import (
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/color"
	"github.com/boggydigital/compton/consts/direction"
	"github.com/boggydigital/compton/consts/font_weight"
	"github.com/boggydigital/compton/consts/size"
	"github.com/boggydigital/compton/elements/els"
	"github.com/boggydigital/compton/elements/flex_items"
	"github.com/boggydigital/compton/elements/fspan"
	"github.com/boggydigital/yet_urls/youtube_urls"
)

func Video(r compton.Registrar, videoId string) compton.Element {
	src := "/video?id=" + videoId
	posterSrc := "/thumbnails?id=" + videoId

	stack := flex_items.FlexItems(r, direction.Column)

	video := els.Video(src)
	video.SetAttribute("controls", "")
	video.SetAttribute("poster", posterSrc)
	stack.Append(video)

	originLink := els.A(youtube_urls.VideoUrl(videoId).String())
	originLink.SetAttribute("target", "_top")
	linkText := fspan.Text(r, "Watch at origin").
		FontSize(size.Small).
		FontWeight(font_weight.Bolder).
		ForegroundColor(color.Cyan)
	originLink.Append(linkText)
	stack.Append(originLink)

	return stack
}

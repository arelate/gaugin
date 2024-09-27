package gaugin_styles

import _ "embed"

var (
	//go:embed "style/app.css"
	AppStyle []byte
	//go:embed "style/screenshots.css"
	ScreenshotsStyle []byte
	//go:embed "style/videos.css"
	VideosStyle []byte
	//go:embed "style/description.css"
	DescriptionStyle []byte
	//go:embed "style/changelog.css"
	ChangelogStyle []byte
	//go:embed "style/steam-deck.css"
	SteamDeckStyle []byte
	//go:embed "style/steam-reviews.css"
	SteamReviews []byte
	//go:embed "style/steam-news.css"
	SteamNews []byte
	//go:embed "style/tag-editors.css"
	TagEditors []byte
)

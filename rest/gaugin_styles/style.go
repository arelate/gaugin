package gaugin_styles

import _ "embed"

var (
	//go:embed "style/app.css"
	AppStyle []byte
	//go:embed "style/screenshots.css"
	ScreenshotsStyle []byte
)

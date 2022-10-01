package stencil_app

// National flags that correspond to language code.
// In some cases there is no obvious way to map those,
// attempting to use sensible option: ar, ca, fa
// Few options are not possible to map to countries (left as comments below)
var LanguageFlags = map[string]string{
	"en":    "🇺🇸",
	"de":    "🇩🇪",
	"fr":    "🇫🇷",
	"es":    "🇪🇸",
	"ru":    "🇷🇺",
	"it":    "🇮🇹",
	"cn":    "🇨🇳",
	"jp":    "🇯🇵",
	"pl":    "🇵🇱",
	"br":    "🇧🇷",
	"ko":    "🇰🇷",
	"zh":    "🇨🇳",
	"tr":    "🇹🇷",
	"cz":    "🇨🇿",
	"pt":    "🇵🇹",
	"nl":    "🇳🇱",
	"es_mx": "🇲🇽",
	"hu":    "🇭🇺",
	"uk":    "🇺🇦",
	"ar":    "🇸🇦",
	"sv":    "🇸🇪",
	"no":    "🇳🇴",
	"da":    "🇩🇰",
	"fi":    "🇫🇮",
	"th":    "🇹🇭",
	"ro":    "🇷🇴",
	"gk":    "🇬🇷",
	"bl":    "🇧🇬",
	"sk":    "🇸🇮",
	"be":    "🇧🇾",
	"he":    "🇮🇱",
	"sb":    "🇷🇸",
	"ca":    "🇪🇸",
	"is":    "🇮🇸",
	"fa":    "🇮🇷",
	"et":    "🇪🇪",
	//Inuktitut (gog_IN): 2 items
	//latine (la): 1 items
}

func LanguageCodeFlag(lc string) string {
	if flag, ok := LanguageFlags[lc]; ok {
		return flag
	} else {
		return lc
	}
}

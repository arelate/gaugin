package api

import (
	"net/http"
)

func HandleFuncs() {

	patternHandlers := map[string]func(w http.ResponseWriter, r *http.Request){
		"/account":        GetAccount,
		"/store":          GetStore,
		"/product":        GetProduct,
		"/search":         GetSearch,
		"/css/styles.css": GetFile,
	}

	for p, h := range patternHandlers {
		http.HandleFunc(p, h)
	}
}

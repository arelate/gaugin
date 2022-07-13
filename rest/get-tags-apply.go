package rest

import (
	"io"
	"net/http"
)

func PostTagsApply(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "applying tags...")
}

package rest

import (
	"fmt"
	"github.com/boggydigital/nod"
	"net/http"
)

func GetMethodOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			err := fmt.Errorf("unsupported method")
			http.Error(w, nod.Error(err).Error(), http.StatusMethodNotAllowed)
			return
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

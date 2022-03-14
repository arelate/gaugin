package api

import (
	"crypto/sha256"
	"crypto/subtle"
	"net/http"
)

func basicHttpAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if u, p, ok := r.BasicAuth(); ok {
			uh, ph := sha256.Sum256([]byte(u)), sha256.Sum256([]byte(p))

			um := subtle.ConstantTimeCompare(uh[:], usernameHash[:]) == 1
			pm := subtle.ConstantTimeCompare(ph[:], passwordHash[:]) == 1

			if um && pm {
				next.ServeHTTP(w, r)
				return
			}
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}

package api

import (
	"bytes"
	"net/http"
	"time"
)

type bytesWriter struct {
	bytes []byte
}

func (bw *bytesWriter) Header() http.Header {
	return http.Header{}
}

func (bw *bytesWriter) WriteHeader(int) {

}

func (bw *bytesWriter) Write(b []byte) (int, error) {
	bw.bytes = append(bw.bytes, b...)
	return len(b), nil
}

var urlStaticCache = make(map[string]*bytesWriter)

func memCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		u := r.URL.String()
		if payload, ok := urlStaticCache[u]; ok {
			http.ServeContent(w, r, u, time.Now(), bytes.NewReader(payload.bytes))
		} else {
			urlStaticCache[u] = &bytesWriter{bytes: make([]byte, 0, 8*1024*1024)}
			next.ServeHTTP(urlStaticCache[u], r)
			http.ServeContent(w, r, u, time.Now(), bytes.NewReader(urlStaticCache[u].bytes))
		}
	})
}

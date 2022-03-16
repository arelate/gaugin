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

var (
	urlStaticCache = make(map[string]*bytesWriter)
	urlTimestamp   = make(map[string]time.Time)
)

func memCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		u := r.URL.String()
		if _, ok := urlStaticCache[u]; !ok {
			urlStaticCache[u] = &bytesWriter{bytes: make([]byte, 0, 1024*1024)}
			urlTimestamp[u] = time.Now()
			next.ServeHTTP(urlStaticCache[u], r)
		}
		w.Header().Set("Cache-Control", "max-age=604800")
		http.ServeContent(w, r, u, urlTimestamp[u], bytes.NewReader(urlStaticCache[u].bytes))
	})
}

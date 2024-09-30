package middleware

import (
	"net/http"
	"strings"
)

const authHeader = "Authorization"

type Base struct {
	authKey string
}

func NewBase(authKey string) *Base {
	return &Base{
		authKey: authKey,
	}
}

func (b *Base) Auth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		header := r.Header.Get(authHeader)
		if header == "" {
			http.Error(w, `{"error": "missing auth key"}`, http.StatusUnauthorized)
			return
		}
		token := strings.Split(header, " ")[1]
		if token != b.authKey {
			http.Error(w, `{"error": "invalid auth key"}`, http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r)
	}
}

// func (b *Base) Headers(h http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Add("Strict-Transport-Security", "max-age=63072000")
// 		h.ServeHTTP(w, r)
// 	})
// }

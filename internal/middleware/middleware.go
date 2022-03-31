package middleware

import (
	"net/http"

	"github.com/KaiserWerk/envy/internal/configuration"
	"github.com/KaiserWerk/envy/internal/logging"
)

type Base struct {
	Config *configuration.AppConfig
	Logger *logging.Logger
}

func (b *Base) Auth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: implement
		h.ServeHTTP(w, r)
	})
}
func (b *Base) Headers(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Strict-Transport-Security", "max-age=63072000")
		h.ServeHTTP(w, r)
	})
}

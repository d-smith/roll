package authzwrapper

import (
	"net/http"
)

type AuthHandler struct {
	handler http.Handler
}

func Wrap(h http.Handler) *AuthHandler {
	return &AuthHandler {handler:h}
}

func (ah *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	authzHeader := r.Header.Get("Authorization")
	if authzHeader ==  "" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized\n"))
		return
	}

	ah.handler.ServeHTTP(w,r)

}
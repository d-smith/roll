package http

import (
	"errors"
	"github.com/xtraclabs/roll/roll"
	"net/http"
)

func handleDevelopers(core *roll.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleDevelopersGet(core, w, r)
		default:
			respondError(w, http.StatusMethodNotAllowed, errors.New("Method not allowed"))
		}
	})
}

func handleDevelopersGet(core *roll.Core, w http.ResponseWriter, r *http.Request) {
	respondOk(w, nil)
}

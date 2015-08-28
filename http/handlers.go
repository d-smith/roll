package http

import (
	"github.com/xtraclabs/roll/roll"
	"net/http"
)

func Handler(core *roll.Core) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/v1/developers", handleDevelopers(core))
	return mux
}

package http

import (
	"github.com/xtraclabs/roll/roll"
	"net/http"
)

func Handler(core *roll.Core) http.Handler {
	mux := http.NewServeMux()
	mux.Handle(DevelopersBaseUri, handleDevelopers(core))
	mux.Handle(ApplicationsBaseUri, handleApplications(core))
	return mux
}

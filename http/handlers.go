package http

import (
	"github.com/xtraclabs/roll/roll"
	"net/http"
)

//Handler creates a much with handlers for all routes in the roll application
func Handler(core *roll.Core) http.Handler {
	mux := http.NewServeMux()
	mux.Handle(DevelopersBaseURI, handleDevelopersBase(core))
	mux.Handle(DevelopersURI, handleDevelopers(core))
	mux.Handle(ApplicationsURI, handleApplications(core))
	mux.Handle(ApplicationsBaseURI, handleApplicationsBase(core))
	mux.Handle(AuthorizeBaseURI, handleAuthorize(core))
	mux.Handle(ValidateBaseURI, handleValidate(core))
	mux.Handle(OAuth2TokenBaseURI, handleToken(core))
	mux.Handle(JWTFlowCertsURI, handleJWTFlowCerts(core))
	mux.Handle(TokenInfoURI, handleTokenInfo(core))
	return mux
}

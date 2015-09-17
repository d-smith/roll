package http

import (
	"github.com/xtraclabs/roll/roll"
	"net/http"
)

//Handler creates a much with handlers for all routes in the roll application
func Handler(core *roll.Core) http.Handler {
	mux := http.NewServeMux()
	mux.Handle(DevelopersBaseURI, handleDevelopers(core))
	mux.Handle(ApplicationsBaseURI, handleApplications(core))
	mux.Handle(AuthorizeBaseURI, handleAuthorize(core))
	mux.Handle(ValidateBaseURI, handleValidate(core))
	mux.Handle(OAuth2TokenBaseURI, handleToken(core))
	mux.Handle(JWTFlowCertsURI, handleJWTFlowCerts(core))
	mux.Handle(TokenInfoURI, handleTokenInfo(core))
	return mux
}

package http

import (
	"github.com/xtraclabs/roll/authzwrapper"
	"github.com/xtraclabs/roll/roll"
	"net/http"
)

//Handler creates a much with handlers for all routes in the roll application
func Handler(core *roll.Core) http.Handler {
	mux := http.NewServeMux()

	//Wrap roll services with the auth checker if booted in secure mode
	if core.Secure() {
		mux.Handle(DevelopersBaseURI, authzwrapper.Wrap(core.SecretsRepo, handleDevelopersBase(core)))
		mux.Handle(DevelopersURI, authzwrapper.Wrap(core.SecretsRepo, handleDevelopers(core)))
		mux.Handle(ApplicationsURI, authzwrapper.Wrap(core.SecretsRepo, handleApplications(core)))
		mux.Handle(ApplicationsBaseURI, authzwrapper.Wrap(core.SecretsRepo, handleApplicationsBase(core)))
	} else {
		mux.Handle(DevelopersBaseURI, handleDevelopersBase(core))
		mux.Handle(DevelopersURI, handleDevelopers(core))
		mux.Handle(ApplicationsURI, handleApplications(core))
		mux.Handle(ApplicationsBaseURI, handleApplicationsBase(core))
	}

	mux.Handle(AuthorizeBaseURI, handleAuthorize(core))
	mux.Handle(ValidateBaseURI, handleValidate(core))
	mux.Handle(OAuth2TokenBaseURI, handleToken(core))
	mux.Handle(JWTFlowCertsURI, handleJWTFlowCerts(core))
	mux.Handle(TokenInfoURI, handleTokenInfo(core))
	return mux
}

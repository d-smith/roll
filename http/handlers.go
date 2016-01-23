package http

import (
	"errors"
	"github.com/xtraclabs/roll/authzwrapper"
	"github.com/xtraclabs/roll/roll"
	"net/http"
	"os"
)

//Handler creates a much with handlers for all routes in the roll application
func Handler(core *roll.Core) http.Handler {
	mux := http.NewServeMux()

	//Wrap roll services with the auth checker if booted in secure mode
	if core.Secure() {
		rollClientID := os.Getenv("ROLL_CLIENTID")
		if rollClientID == "" {
			panic(errors.New("Cannot run in secure mode without a client ID to white list (from ROLL_CLIENTID env variable)"))
		}

		whitelist := []string{rollClientID}
		mux.Handle(DevelopersBaseURI, authzwrapper.Wrap(core.SecretsRepo, core.AdminRepo, whitelist, handleDevelopersBase(core)))
		mux.Handle(DevelopersURI, authzwrapper.Wrap(core.SecretsRepo, core.AdminRepo, whitelist, handleDevelopers(core)))
		mux.Handle(ApplicationsURI, authzwrapper.Wrap(core.SecretsRepo, core.AdminRepo, whitelist, handleApplications(core)))
		mux.Handle(ApplicationsBaseURI, authzwrapper.Wrap(core.SecretsRepo, core.AdminRepo, whitelist, handleApplicationsBase(core)))
		mux.Handle(JWTFlowCertsURI, authzwrapper.Wrap(core.SecretsRepo, core.AdminRepo, whitelist, handleJWTFlowCerts(core)))
	} else {
		mux.Handle(DevelopersBaseURI, authzwrapper.WrapUnsecure(handleDevelopersBase(core)))
		mux.Handle(DevelopersURI, authzwrapper.WrapUnsecure(handleDevelopers(core)))
		mux.Handle(ApplicationsURI, authzwrapper.WrapUnsecure(handleApplications(core)))
		mux.Handle(ApplicationsBaseURI, authzwrapper.WrapUnsecure(handleApplicationsBase(core)))
		mux.Handle(JWTFlowCertsURI, authzwrapper.WrapUnsecure(handleJWTFlowCerts(core)))
	}

	mux.Handle(AuthorizeBaseURI, handleAuthorize(core))
	mux.Handle(ValidateBaseURI, handleValidate(core))
	mux.Handle(OAuth2TokenBaseURI, handleToken(core))
	mux.Handle(TokenInfoURI, handleTokenInfo(core))
	return mux
}

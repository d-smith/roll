package authzwrapper

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/context"
	"net/http"
)

const UnsecureRollSubjectHeader = "X-Roll-Subject"

type unsecureHandler struct {
	handler http.Handler
}

func WrapUnsecure(h http.Handler) http.Handler {

	return &unsecureHandler{
		handler: h,
	}
}

func (uh unsecureHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	log.Info("unsecure wrapped")
	subject := r.Header.Get(UnsecureRollSubjectHeader)
	if subject == "" {
		log.Info("Missing Authorization header (unsecure mode)")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized\n"))
		return
	}

	log.Info("setting subject and context:", subject, " ", false)
	context.Set(r, AuthzSubject, subject)
	context.Set(r, AuthzAdminScope, false)
	uh.handler.ServeHTTP(w, r)

	log.Info("clear context")
	context.Clear(r)
}

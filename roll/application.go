package roll

//Application represents the data associated with an application that is exposed via the REST API
type Application struct {
	DeveloperEmail   string
	CLientID         string
	ApplicationName  string
	ClientSecret     string
	RedirectURI      string
	LoginProvider    string
	JWTFlowPublicKey string
}

//ApplicationRepo represents a repository abstraction for dealing with persistent Application instances.
type ApplicationRepo interface {
	StoreApplication(app *Application) error
	RetrieveApplication(clientID string) (*Application, error)
}

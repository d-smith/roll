package roll

//Application represents the data associated with an application that is exposed via the REST API
type Application struct {
	DeveloperEmail   string
	APIKey           string
	ApplicationName  string
	APISecret        string
	RedirectURI      string
	LoginProvider    string
	JWTFlowPublicKey string
}

//ApplicationRepo represents a repository abstraction for dealing with persistent Application instances.
type ApplicationRepo interface {
	StoreApplication(app *Application) error
	RetrieveApplication(apiKey string) (*Application, error)
}

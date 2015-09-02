package roll

type Application struct {
	DeveloperEmail  string
	APIKey          string
	ApplicationName string
	APISecret       string
	RedirectUri		string
}

type ApplicationRepo interface {
	StoreApplication(app *Application) error
	RetrieveApplication(apiKey string) (*Application, error)
}

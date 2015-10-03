package repos

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/xtraclabs/roll/dbutil"
	"github.com/xtraclabs/roll/roll"
	"github.com/xtraclabs/roll/secrets"
	"log"
)

//DynamoAppRepo presents a repository interface for storing and retrieving application definitions,
//backed by DynamoDB
type DynamoAppRepo struct {
	client *dynamodb.DynamoDB
}

//NewDynamoAppRepo returns a new instance of type DynamoAppRepo
func NewDynamoAppRepo() *DynamoAppRepo {
	return &DynamoAppRepo{
		client: dbutil.CreateDynamoDBClient(),
	}
}

//StoreApplication stores an application definition in DynamoDB
func (dar *DynamoAppRepo) StoreApplication(app *roll.Application) error {

	if app.ClientSecret == "" {
		clientSecret, err := secrets.GenerateClientSecret()
		if err != nil {
			return err
		}
		app.ClientSecret = clientSecret
	}

	appAttrs := map[string]*dynamodb.AttributeValue{
		"ClientID":        {S: aws.String(app.ClientID)},
		"ApplicationName": {S: aws.String(app.ApplicationName)},
		"ClientSecret":    {S: aws.String(app.ClientSecret)},
		"DeveloperEmail":  {S: aws.String(app.DeveloperEmail)},
		"RedirectUri":     {S: aws.String(app.RedirectURI)},
		"LoginProvider":   {S: aws.String(app.LoginProvider)},
	}

	if app.JWTFlowPublicKey != "" {
		appAttrs["JWTFlowPublicKey"] = &dynamodb.AttributeValue{
			S: aws.String(app.JWTFlowPublicKey),
		}
	}

	params := &dynamodb.PutItemInput{
		TableName: aws.String("Application"),
		Item:      appAttrs,
	}
	_, err := dar.client.PutItem(params)

	return err
}

//RetrieveApplication retrieves an application definition from DynamoDB
func (dar *DynamoAppRepo) RetrieveApplication(clientID string) (*roll.Application, error) {
	params := &dynamodb.GetItemInput{
		TableName: aws.String("Application"),
		Key: map[string]*dynamodb.AttributeValue{
			"ClientID": {S: aws.String(clientID)},
		},
	}

	log.Println("Get item")
	out, err := dar.client.GetItem(params)
	if err != nil {
		return nil, err
	}

	if len(out.Item) == 0 {
		return nil, nil
	}

	log.Println("Load struct with data returned from dynamo")
	return &roll.Application{
		ClientID:         extractString(out.Item["ClientID"]),
		ApplicationName:  extractString(out.Item["ApplicationName"]),
		ClientSecret:     extractString(out.Item["ClientSecret"]),
		DeveloperEmail:   extractString(out.Item["DeveloperEmail"]),
		RedirectURI:      extractString(out.Item["RedirectUri"]),
		LoginProvider:    extractString(out.Item["LoginProvider"]),
		JWTFlowPublicKey: extractString(out.Item["JWTFlowPublicKey"]),
	}, nil
}

func (dar *DynamoAppRepo) ListApplications() ([]roll.Application, error) {
	params := &dynamodb.ScanInput{
		TableName: aws.String("Application"),
		AttributesToGet: []*string{
			aws.String("ClientID"),
			aws.String("ApplicationName"),
			aws.String("ClientSecret"),
			aws.String("DeveloperEmail"),
			aws.String("RedirectUri"),
			aws.String("LoginProvider"),
			aws.String("JWTFlowPublicKey"),
		},
	}

	resp, err := dar.client.Scan(params)
	if err != nil {
		return nil, err
	}

	var apps []roll.Application

	for _, item := range resp.Items {
		application := roll.Application{
			ClientID:         extractString(item["ClientID"]),
			ApplicationName:  extractString(item["ApplicationName"]),
			ClientSecret:     extractString(item["ClientSecret"]),
			DeveloperEmail:   extractString(item["DeveloperEmail"]),
			RedirectURI:      extractString(item["RedirectUri"]),
			LoginProvider:    extractString(item["LoginProvider"]),
			JWTFlowPublicKey: extractString(item["JWTFlowPublicKey"]),
		}

		apps = append(apps, application)
	}
	return apps, nil
}

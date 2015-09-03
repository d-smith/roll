package repos

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
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
	//TODO - pick up region from config?
	return &DynamoAppRepo{
		client: dynamodb.New(&aws.Config{Region: aws.String("us-east-1")}),
	}
}

//StoreApplication stores an application definition in DynamoDB
func (dar *DynamoAppRepo) StoreApplication(app *roll.Application) error {
	//TODO - do we generate a secret everytime this is called? Probably need a POST to
	//create and a put to update - refactor later after talking this through with others

	apiSecret, err := secrets.GenerateApiSecret()
	if err != nil {
		return err
	}

	params := &dynamodb.PutItemInput{
		TableName: aws.String("Application"),
		Item: map[string]*dynamodb.AttributeValue{
			"APIKey":          {S: aws.String(app.APIKey)},
			"ApplicationName": {S: aws.String(app.ApplicationName)},
			"APISecret":       {S: aws.String(apiSecret)},
			"DeveloperEmail":  {S: aws.String(app.DeveloperEmail)},
			"RedirectUri":     {S: aws.String(app.RedirectURI)},
		},
	}
	_, err = dar.client.PutItem(params)

	return err
}

//RetrieveApplication retrieves an application definition from DynamoDB
func (dar *DynamoAppRepo) RetrieveApplication(apiKey string) (*roll.Application, error) {
	params := &dynamodb.GetItemInput{
		TableName: aws.String("Application"),
		Key: map[string]*dynamodb.AttributeValue{
			"APIKey": {S: aws.String(apiKey)},
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
		APIKey:          extractString(out.Item["APIKey"]),
		ApplicationName: extractString(out.Item["ApplicationName"]),
		APISecret:       extractString(out.Item["APISecret"]),
		DeveloperEmail:  extractString(out.Item["DeveloperEmail"]),
		RedirectURI:     extractString(out.Item["RedirectUri"]),
	}, nil
}

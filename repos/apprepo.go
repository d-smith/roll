package repos

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/xtraclabs/roll/roll"
	"errors"
	"github.com/xtraclabs/roll/secrets"
)

type DynamoAppRepo struct {
	client *dynamodb.DynamoDB
}

func NewDynamoAppRepo() *DynamoAppRepo {
	//TODO - pick up region from config?
	return &DynamoAppRepo{
		client: dynamodb.New(&aws.Config{Region: aws.String("us-east-1")}),
	}
}

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
			"APIKey":     {S: aws.String(app.APIKey)},
			"ApplicationName":     {S: aws.String(app.ApplicationName)},
			"APISecret":     {S: aws.String(apiSecret)},
			"DeveloperEmail": {S: aws.String(app.DeveloperEmail)},
		},
	}
	_, err = dar.client.PutItem(params)

	return err
}

func (dar *DynamoAppRepo)RetrieveApplication(apiKey string) (*roll.Application, error) {
	return nil,  errors.New("RetrieveApplication not implemented")
}
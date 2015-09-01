package repos

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/xtraclabs/roll/roll"
	"errors"
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
	return errors.New("StoreApplication not implemented")
}

func (dar *DynamoAppRepo)RetrieveApplication(apiKey string) (*roll.Application, error) {
	return nil,  errors.New("RetrieveApplication not implemented")
}
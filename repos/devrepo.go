package repos

import (
	"log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/xtraclabs/roll/roll"
)

type DynamoDevRepo struct {
	client *dynamodb.DynamoDB
}

func extractString(attrval *dynamodb.AttributeValue) string {
	if attrval == nil {
		return ""
	}

	return *attrval.S
}

func NewDynamoDevRepo() *DynamoDevRepo {
	//TODO - pick up region from config?
	return &DynamoDevRepo{
		client: dynamodb.New(&aws.Config{Region: aws.String("us-east-1")}),
	}
}

func (dddr DynamoDevRepo) RetrieveDeveloper(email string) (*roll.Developer, error) {
	params := &dynamodb.GetItemInput{
		TableName: aws.String("Developer"),
		Key: map[string]*dynamodb.AttributeValue{
			"EMail": {S: aws.String(email)},
		},
	}

	log.Println("Get item")
	out, err := dddr.client.GetItem(params)
	if err != nil {
		return nil, err
	}

	if len(out.Item) == 0 {
		return nil, nil
	}

	log.Println("Load struct with data returned from dynamo")
	return &roll.Developer{
		Email:     extractString(out.Item["EMail"]),
		FirstName: extractString(out.Item["FirstName"]),
		LastName:  extractString(out.Item["LastName"]),
	}, nil
}


func (dddr DynamoDevRepo) StoreDeveloper(dev *roll.Developer) error {
	params := &dynamodb.PutItemInput{
		TableName: aws.String("Developer"),
		Item: map[string]*dynamodb.AttributeValue{
			"EMail":     {S: aws.String(dev.Email)},
			"FirstName": {S: aws.String(dev.FirstName)},
			"LastName":  {S: aws.String(dev.LastName)},
		},
	}
	_, err := dddr.client.PutItem(params)

	return err
}
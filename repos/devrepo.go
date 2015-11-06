package repos

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/xtraclabs/roll/dbutil"
	"github.com/xtraclabs/roll/roll"
	"log"
)

//DynamoDevRepo provides a repository for Developer objects implemented using DynamoDB
type DynamoDevRepo struct {
	client *dynamodb.DynamoDB
}

func extractString(attrval *dynamodb.AttributeValue) string {
	if attrval == nil {
		return ""
	}

	return *attrval.S
}

//NewDynamoDevRepo creates a new instance of DynamoDevRepo
func NewDynamoDevRepo() *DynamoDevRepo {
	return &DynamoDevRepo{
		client: dbutil.CreateDynamoDBClient(),
	}
}

const (
	EMail     = "EMail"
	FirstName = "FirstName"
	LastName  = "LastName"
	ID        = "ID"
)

//RetrieveDeveloper retrieves a developer from DynamoDB using the developer's email as the key
func (dddr DynamoDevRepo) RetrieveDeveloper(email string, subjectID string, adminScope bool) (*roll.Developer, error) {
	params := &dynamodb.QueryInput{
		TableName:              aws.String("Developer"),
		KeyConditionExpression: aws.String("EMail=:email"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":email": {S: aws.String(email)},
		},
		IndexName: aws.String("Email-Index"),
	}

	if !adminScope {
		params.FilterExpression = aws.String("ID=:subjectID")
		params.ExpressionAttributeValues = map[string]*dynamodb.AttributeValue{
			":email":     {S: aws.String(email)},
			":subjectID": {S: aws.String(subjectID)},
		}
	}

	resp, err := dddr.client.Query(params)
	if err != nil {
		return nil, err
	}

	if *resp.Count == 0 {
		return nil, nil
	}

	if *resp.Count > 1 {
		return nil, fmt.Errorf("Expected 1 result got %d instead", *resp.Count)
	}

	log.Println("Load struct with data returned from dynamo")
	return &roll.Developer{
		Email:     extractString(resp.Items[0][EMail]),
		FirstName: extractString(resp.Items[0][FirstName]),
		LastName:  extractString(resp.Items[0][LastName]),
		ID:        extractString(resp.Items[0][ID]),
	}, nil
}

//StoreDeveloper stores a developer instance in dynamoDB
func (dddr DynamoDevRepo) StoreDeveloper(dev *roll.Developer) error {
	params := &dynamodb.PutItemInput{
		TableName: aws.String("Developer"),
		Item: map[string]*dynamodb.AttributeValue{
			EMail:     {S: aws.String(dev.Email)},
			FirstName: {S: aws.String(dev.FirstName)},
			LastName:  {S: aws.String(dev.LastName)},
			ID:        {S: aws.String(dev.ID)},
		},
	}
	_, err := dddr.client.PutItem(params)

	return err
}

func (dddr DynamoDevRepo) ListDevelopers(subjectID string, adminScope bool) ([]roll.Developer, error) {
	params := &dynamodb.ScanInput{
		TableName: aws.String("Developer"),
	}

	if !adminScope {
		params.FilterExpression = aws.String("ID=:subjectID")
		params.ExpressionAttributeValues = map[string]*dynamodb.AttributeValue{
			":subjectID": {S: aws.String(subjectID)},
		}
	}

	resp, err := dddr.client.Scan(params)
	if err != nil {
		return nil, err
	}

	var devs []roll.Developer

	for _, item := range resp.Items {
		developer := roll.Developer{
			Email:     extractString(item[EMail]),
			FirstName: extractString(item[FirstName]),
			LastName:  extractString(item[LastName]),
			ID:        extractString(item[ID]),
		}

		devs = append(devs, developer)
	}
	return devs, nil
}

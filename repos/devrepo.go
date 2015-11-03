package repos

import (
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
		Email:     extractString(out.Item[EMail]),
		FirstName: extractString(out.Item[FirstName]),
		LastName:  extractString(out.Item[LastName]),
		ID:        extractString(out.Item[ID]),
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

func (dddr DynamoDevRepo) ListDevelopers() ([]roll.Developer, error) {
	params := &dynamodb.ScanInput{
		TableName: aws.String("Developer"),
		AttributesToGet: []*string{
			aws.String(EMail),
			aws.String(FirstName),
			aws.String(LastName),
			aws.String(ID),
		},
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

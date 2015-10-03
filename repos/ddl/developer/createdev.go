package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"log"
	"os"
)

func main() {
	var svc *dynamodb.DynamoDB
	localAddr := os.Getenv("LOCAL_DYNAMO_ADDR")
	if localAddr != "" {
		svc = dynamodb.New(&aws.Config{Endpoint: aws.String(localAddr), Region: aws.String("here")})
	} else {
		svc = dynamodb.New(&aws.Config{Region: aws.String("us-east-1")})
	}

	params := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("EMail"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("EMail"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
		TableName: aws.String("Developer"),
	}

	resp, err := svc.CreateTable(params)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(resp)
}

package dbutil

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"log"
	"os"
)

func CreateDynamoDBClient() *dynamodb.DynamoDB {
	var dynamoClient *dynamodb.DynamoDB


	localAddr := os.Getenv("LOCAL_DYNAMO_ADDR") //e.g. http://locahost:8000
	if localAddr != "" {
		log.Printf("Using local dynamodb address - %s", localAddr)
		dynamoClient = dynamodb.New(&aws.Config{Endpoint: aws.String(localAddr), Region: aws.String("here")})
	} else {
		dynamoClient = dynamodb.New(&aws.Config{Region: aws.String("us-east-1")})
	}

	return dynamoClient
}

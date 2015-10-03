package dbutil

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"os"
)

func CreateDynamoDBClient() *dynamodb.DynamoDB {
	//TODO - what's the best way to pick up AWS configuration?
	var dynamoClient *dynamodb.DynamoDB

	localAddr := os.Getenv("LOCAL_DYNAMO_ADDR")
	if localAddr != "" {
		dynamoClient = dynamodb.New(&aws.Config{Endpoint: aws.String(localAddr), Region: aws.String("here")})
	} else {
		dynamoClient = dynamodb.New(&aws.Config{Region: aws.String("us-east-1")})
	}

	return dynamoClient
}

package dbutil

import (
	"database/sql"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

func CreateDynamoDBClient() *dynamodb.DynamoDB {
	var dynamoClient *dynamodb.DynamoDB

	localAddr := os.Getenv("LOCAL_DYNAMO_ADDR") //e.g. export LOCAL_DYNAMO_ADDR=http://localhost:8000
	if localAddr != "" {
		log.Printf("Using local dynamodb address - %s", localAddr)
		dynamoClient = dynamodb.New(&aws.Config{Endpoint: aws.String(localAddr), Region: aws.String("here")})
	} else {
		dynamoClient = dynamodb.New(&aws.Config{Region: aws.String("us-east-1")})
	}

	return dynamoClient
}

func CreateMariaDBSqlDB() (*sql.DB, error) {
	rollUser := os.Getenv("ROLL_DBUSER")
	rollPassword := os.Getenv("ROLL_DBPASSWORD")
	rollAddress := os.Getenv("ROLL_ADDRESS")

	connectString := fmt.Sprintf("%s:%s@tcp(%s)/rolldb", rollUser, rollPassword, rollAddress)

	return sql.Open("mysql", connectString)
}

package ddl

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/xtraclabs/roll/dbutil"
	"log"
	"os"
)

const (
	//DynamoDB table name for holding admins
	AdminTableName = "Admin"

	//DynamoDB table name for storing registered application details
	ApplicationTableName = "Application"

	//DynamoDB table name for storing registered developers
	DeveloperTableName = "Developer"
)

func DeleteTable(tableName string) {

	localAddr := os.Getenv("LOCAL_DYNAMO_ADDR")
	if localAddr == "" {
		log.Println("DeleteAppTable will only attempt to delete a local dynamodb table... returning.")
		return
	}

	var svc *dynamodb.DynamoDB = dbutil.CreateDynamoDBClient()

	params := &dynamodb.DeleteTableInput{
		TableName: aws.String(tableName),
	}
	resp, err := svc.DeleteTable(params)

	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Printf("Table %s deleted:\n%s\n", tableName, resp.String())
}

func CreateAppTable() {
	var svc *dynamodb.DynamoDB = dbutil.CreateDynamoDBClient()

	params := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("ClientID"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("DeveloperEmail"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("ClientID"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
		TableName: aws.String(ApplicationTableName),
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			{ // Required
				IndexName: aws.String("EMail-Index"),
				KeySchema: []*dynamodb.KeySchemaElement{
					{ // Required
						AttributeName: aws.String("DeveloperEmail"),
						KeyType:       aws.String("HASH"),
					},
					// More values...
				},
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("ALL"),
				},
				ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(1),
					WriteCapacityUnits: aws.Int64(1),
				},
			},
		},
	}

	resp, err := svc.CreateTable(params)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(resp)
}

func CreateDevTable() {
	var svc *dynamodb.DynamoDB = dbutil.CreateDynamoDBClient()

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
		TableName: aws.String(DeveloperTableName),
	}

	resp, err := svc.CreateTable(params)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(resp)
}

func CreateAdminTable() {
	var svc *dynamodb.DynamoDB = dbutil.CreateDynamoDBClient()

	params := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("AdminID"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("AdminID"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
		TableName: aws.String(AdminTableName),
	}

	resp, err := svc.CreateTable(params)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(resp)
}

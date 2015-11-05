package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/xtraclabs/roll/dbutil"
	"github.com/xtraclabs/roll/repos/ddl"
	"os"
)

func main() {
	var add = flag.String("add", "", "Add a subject as an admin")
	var remove = flag.String("remove", "", "Remove a subject as an admin")
	var list = flag.Bool("list", false, "List admins")
	flag.Parse()

	doList := ""
	if *list {
		doList = "dolist"
	}

	if !oneOperationRequested([]string{*add, *remove, doList}) {
		fmt.Println("Specify a single action please")
		os.Exit(1)
	}

	if *add != "" {
		handleAdd(*add)
	} else if *remove != "" {
		handleRemove(*remove)
	} else if *list {
		handleList()
	}

	os.Exit(0)
}

func oneOperationRequested(ops []string) bool {
	var count = 0
	for _, op := range ops {
		if op != "" {
			count += 1
		}
	}

	return count == 1
}

func handleAdd(subject string) {
	fmt.Println("Add admin", subject)

	var client *dynamodb.DynamoDB = dbutil.CreateDynamoDBClient()
	params := &dynamodb.PutItemInput{
		TableName: aws.String(ddl.AdminTableName),
		Item: map[string]*dynamodb.AttributeValue{
			"AdminID": {S: aws.String(subject)},
		},
		Expected: map[string]*dynamodb.ExpectedAttributeValue{
			"AdminID": {Exists: aws.Bool(false)},
		},
	}
	_, err := client.PutItem(params)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func handleRemove(subject string) {
	fmt.Println("Remove admin", subject)

	var client *dynamodb.DynamoDB = dbutil.CreateDynamoDBClient()
	params := &dynamodb.DeleteItemInput{
		TableName: aws.String(ddl.AdminTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"AdminID": {S: aws.String(subject)},
		},
	}

	_, err := client.DeleteItem(params)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func handleList() {
	var client *dynamodb.DynamoDB = dbutil.CreateDynamoDBClient()
	params := &dynamodb.ScanInput{
		TableName: aws.String(ddl.AdminTableName),
		AttributesToGet: []*string{
			aws.String("AdminID"),
		},
	}

	resp, err := client.Scan(params)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(resp)
}

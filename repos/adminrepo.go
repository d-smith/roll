package repos

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/xtraclabs/roll/dbutil"
	"github.com/xtraclabs/roll/repos/ddl"
)

//DynamoAdminRepo presents a repository interface reading the admin table in DynamoDB to determine if
//a subject can be granted admin scope
type DynamoAdminRepo struct {
	client *dynamodb.DynamoDB
}

//NewDynamoAdminRepo returns a new instance of type DynamoAppRepo
func NewDynamoAdminRepo() *DynamoAdminRepo {
	return &DynamoAdminRepo{
		client: dbutil.CreateDynamoDBClient(),
	}
}

//IsAdmin returns true is the given subject is present in the admin table, and false otherwise
func (ar *DynamoAdminRepo) IsAdmin(subject string) (bool, error) {
	params := &dynamodb.GetItemInput{
		TableName: aws.String(ddl.AdminTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"AdminID": {S: aws.String(subject)},
		},
	}

	out, err := ar.client.GetItem(params)
	if err != nil {
		return false, err
	}

	return len(out.Item) == 1, nil
}

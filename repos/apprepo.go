package repos

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/xtraclabs/roll/dbutil"
	"github.com/xtraclabs/roll/roll"
	"github.com/xtraclabs/roll/secrets"
	"log"
)

//DynamoAppRepo presents a repository interface for storing and retrieving application definitions,
//backed by DynamoDB
type DynamoAppRepo struct {
	client *dynamodb.DynamoDB
}

//NewDynamoAppRepo returns a new instance of type DynamoAppRepo
func NewDynamoAppRepo() *DynamoAppRepo {
	return &DynamoAppRepo{
		client: dbutil.CreateDynamoDBClient(),
	}
}

type DuplicateAppdefError struct {
	ApplicationName string
	DeveloperEmail  string
}

func NewDuplicationAppdefError(appName, devEmail string) *DuplicateAppdefError {
	return &DuplicateAppdefError{
		ApplicationName: appName,
		DeveloperEmail:  devEmail,
	}
}

const (
	ClientID = "ClientID"
	ApplicationName = "ApplicationName"
	ClientSecret = "ClientSecret"
	DeveloperEmail = "DeveloperEmail"
	RedirectUri = "RedirectUri"
	LoginProvider = "LoginProvider"
	JWTFlowPublicKey = "JWTFlowPublicKey"
)

func (dae *DuplicateAppdefError) Error() string {
	return fmt.Sprintf("Application definition exists for application name %s and developer email %s",
		dae.ApplicationName, dae.DeveloperEmail)
}

//CreateApplication stores an application definition in DynamoDB
func (dar *DynamoAppRepo) CreateApplication(app *roll.Application) error {

	//Make sure we are not creating a new application definition for an existing
	//application name/developer email combination
	existing, err := dar.RetrieveAppByNameAndDevEmail(app.ApplicationName, app.DeveloperEmail)
	if err != nil {
		return err
	}

	if existing != nil {
		return NewDuplicationAppdefError(app.ApplicationName, app.DeveloperEmail)
	}

	if app.ClientSecret == "" {
		clientSecret, err := secrets.GenerateClientSecret()
		if err != nil {
			return err
		}
		app.ClientSecret = clientSecret
	}

	appAttrs := map[string]*dynamodb.AttributeValue{
		ClientID:        {S: aws.String(app.ClientID)},
		ApplicationName: {S: aws.String(app.ApplicationName)},
		ClientSecret:    {S: aws.String(app.ClientSecret)},
		DeveloperEmail:  {S: aws.String(app.DeveloperEmail)},
		RedirectUri:     {S: aws.String(app.RedirectURI)},
		LoginProvider:   {S: aws.String(app.LoginProvider)},
	}

	if app.JWTFlowPublicKey != "" {
		appAttrs[JWTFlowPublicKey] = &dynamodb.AttributeValue{
			S: aws.String(app.JWTFlowPublicKey),
		}
	}

	params := &dynamodb.PutItemInput{
		TableName:           aws.String("Application"),
		ConditionExpression: aws.String("attribute_not_exists(ClientID)"),
		Item:                appAttrs,
	}
	_, err = dar.client.PutItem(params)

	return err
}

//UpdateApplication updates an existing application definition
func (dar *DynamoAppRepo) UpdateApplication(app *roll.Application) error {

	//Build up the non-empty attributes to update
	updateAttributes := make(map[string]*dynamodb.AttributeValueUpdate)

	if app.LoginProvider != "" {
		updateAttributes[LoginProvider] = &dynamodb.AttributeValueUpdate{
			Action: aws.String(dynamodb.AttributeActionPut),
			Value: &dynamodb.AttributeValue{
				S: aws.String(app.LoginProvider),
			},
		}
	}

	if app.RedirectURI != "" {
		log.Println("Updating redirect uri:", app.RedirectURI)
		updateAttributes[RedirectUri] = &dynamodb.AttributeValueUpdate{
			Action: aws.String(dynamodb.AttributeActionPut),
			Value: &dynamodb.AttributeValue{
				S: aws.String(app.RedirectURI),
			},
		}
	}

	if app.JWTFlowPublicKey != "" {
		updateAttributes[JWTFlowPublicKey] = &dynamodb.AttributeValueUpdate{
			Action: aws.String(dynamodb.AttributeActionPut),
			Value: &dynamodb.AttributeValue{
				S: aws.String(app.JWTFlowPublicKey),
			},
		}
	}

	if app.ApplicationName != "" {
		log.Println("Updating application name:", app.ApplicationName)
		updateAttributes[ApplicationName] = &dynamodb.AttributeValueUpdate{
			Action: aws.String(dynamodb.AttributeActionPut),
			Value: &dynamodb.AttributeValue{
				S: aws.String(app.ApplicationName),
			},
		}
	}

	params := &dynamodb.UpdateItemInput{
		TableName: aws.String("Application"),
		Key: map[string]*dynamodb.AttributeValue{
			ClientID: {S: aws.String(app.ClientID)},
		},
		AttributeUpdates: updateAttributes,
	}

	_, err := dar.client.UpdateItem(params)

	return err
}

//RetrieveAppByNameAndDevEmail retrieves an application definition based on the combination of
//application name and developer email
func (dar *DynamoAppRepo) RetrieveAppByNameAndDevEmail(appName, email string) (*roll.Application, error) {
	params := &dynamodb.QueryInput{
		TableName:              aws.String("Application"),
		IndexName:              aws.String("EMail-Index"),
		KeyConditionExpression: aws.String("DeveloperEmail=:email"),
		FilterExpression:       aws.String("ApplicationName=:appName"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":email":   {S: aws.String(email)},
			":appName": {S: aws.String(appName)},
		},
	}

	resp, err := dar.client.Query(params)
	if err != nil {
		return nil, err
	}

	if resp == nil || *resp.Count == 0 {
		return nil, nil
	}

	return &roll.Application{
		ClientID:         extractString(resp.Items[0][ClientID]),
		ApplicationName:  extractString(resp.Items[0][ApplicationName]),
		ClientSecret:     extractString(resp.Items[0][ClientSecret]),
		DeveloperEmail:   extractString(resp.Items[0][DeveloperEmail]),
		RedirectURI:      extractString(resp.Items[0][RedirectUri]),
		LoginProvider:    extractString(resp.Items[0][LoginProvider]),
		JWTFlowPublicKey: extractString(resp.Items[0][JWTFlowPublicKey]),
	}, nil
}

//RetrieveApplication retrieves an application definition from DynamoDB. Note a nil
//pointer is returned if a successful call to dynamodb does not find an application
//stored for the given clientID
func (dar *DynamoAppRepo) RetrieveApplication(clientID string) (*roll.Application, error) {
	params := &dynamodb.GetItemInput{
		TableName: aws.String("Application"),
		Key: map[string]*dynamodb.AttributeValue{
			ClientID: {S: aws.String(clientID)},
		},
	}

	log.Println("Get item")
	out, err := dar.client.GetItem(params)
	if err != nil {
		return nil, err
	}

	if len(out.Item) == 0 {
		return nil, nil
	}

	log.Println("Load struct with data returned from dynamo")
	return &roll.Application{
		ClientID:         extractString(out.Item[ClientID]),
		ApplicationName:  extractString(out.Item[ApplicationName]),
		ClientSecret:     extractString(out.Item[ClientSecret]),
		DeveloperEmail:   extractString(out.Item[DeveloperEmail]),
		RedirectURI:      extractString(out.Item[RedirectUri]),
		LoginProvider:    extractString(out.Item[LoginProvider]),
		JWTFlowPublicKey: extractString(out.Item[JWTFlowPublicKey]),
	}, nil
}

func (dar *DynamoAppRepo) ListApplications() ([]roll.Application, error) {
	params := &dynamodb.ScanInput{
		TableName: aws.String("Application"),
		AttributesToGet: []*string{
			aws.String(ClientID),
			aws.String(ApplicationName),
			aws.String(ClientSecret),
			aws.String(DeveloperEmail),
			aws.String(RedirectUri),
			aws.String(LoginProvider),
			aws.String(JWTFlowPublicKey),
		},
	}

	resp, err := dar.client.Scan(params)
	if err != nil {
		return nil, err
	}

	var apps []roll.Application

	for _, item := range resp.Items {
		application := roll.Application{
			ClientID:         extractString(item[ClientID]),
			ApplicationName:  extractString(item[ApplicationName]),
			ClientSecret:     extractString(item[ClientSecret]),
			DeveloperEmail:   extractString(item[DeveloperEmail]),
			RedirectURI:      extractString(item[RedirectUri]),
			LoginProvider:    extractString(item[LoginProvider]),
			JWTFlowPublicKey: extractString(item[JWTFlowPublicKey]),
		}

		apps = append(apps, application)
	}
	return apps, nil
}

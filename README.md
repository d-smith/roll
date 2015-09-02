## Roll - Services to maintain a roll of developers and their applications

Patterns and some code inspired by Hashicorp vault.

### Runtime set up
Configure AWS credentials the normal way - I've run the program with the following user policy:

<pre>
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "Stmt1441122614000",
            "Effect": "Allow",
            "Action": [
                "dynamodb:DeleteItem",
                "dynamodb:GetItem",
                "dynamodb:GetRecords",
                "dynamodb:PutItem",
                "dynamodb:Query",
                "dynamodb:Scan",
                "dynamodb:UpdateItem"
            ],
            "Resource": [
                "*"
            ]
        }
    ]
}
</pre>

This assumes two tables have been created in DynamoDB:

* Developer with a string hash key named EMail
* Application with a string hash key named APIKey

### Build dependencies

Still need to vendor my dependencies, but they are:

* [Stretchr Testify](https://github.com/stretchr/testify/)
* [AWS Golang SDK](https://github.com/aws/aws-sdk-go)

I used the [mockery tool](https://github.com/vektra/mockery) to generate the mocks - I don't think there's a runtime
 dependency but go get it as per its instructions if you see build weirdness.


### Early Look

Register a dev:

<pre>
curl -v -X PUT -d '
{
"Email":"doug@dev.com",
"FirstName":"Doug",
"LastName":"Dev"
}' localhost:12345/v1/developers/doug@dev.com
</pre>

Retrieve a dev:

<pre>
curl -v localhost:12345/v1/developers/doug@dev.com
</pre>


Register an app:

<pre>
curl -X PUT -d '{
"APIKey":"111-222-3333",
"ApplicationName":"Hangry Birdz",
"DeveloperEmail":"doug@dev.com",
"RedirectUri":"http://localhost:3000/ab"
}' localhost:12345/v1/applications/111-222-3333
</pre>

Retrieve registered app:

<pre>
curl -v localhost:12345/v1/applications/111-222-3333
</pre>
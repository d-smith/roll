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
* [Hashicorp Vault Golang API](https://github.com/hashicorp/vault/tree/master/api)
* [JWT-go](https://github.com/dgrijalva/jwt-go)
* [Go UUID](https://github.com/nu7hatch/gouuid)

Use `go get github.com/hashicorp/vault/api` to install the API portion of Vautl

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
}' localhost:3000/v1/developers/doug@dev.com
</pre>

Retrieve a dev:

<pre>
curl -v localhost:3000/v1/developers/doug@dev.com
</pre>


Register an app:

<pre>
curl -X PUT -d '{
"APIKey":"111-222-3333",
"ApplicationName":"Hangry Birdz",
"DeveloperEmail":"doug@dev.com",
"RedirectURI":"http://localhost:2000/oauth2_callback"
}' localhost:3000/v1/applications/111-222-3333
</pre>

Retrieve registered app:

<pre>
curl -v localhost:3000/v1/applications/111-222-3333
</pre>

### Implicit Grant Flow

First, fire up the services - go build in the rollsvcs directory, then run the server

<pre>
./rollsvcs -port 3000
</pre>

You will also need to fire up the callback server - go build in cbserver then run it:

<pre>
./cbserver -port 2000
</pre>

Next register an app, e.g.

<pre>
curl -X PUT -d '{
"APIKey":"111-222-3333",
"ApplicationName":"Hangry Birdz",
"DeveloperEmail":"doug@dev.com",
"RedirectURI":"http://localhost:2000/oauth2_callback",
"LoginProvider":"xtrac://localhost:9000"
}' localhost:3000/v1/applications/111-222-3333
</pre>

For the next step, fire up your browser for point to:

<pre>
http://localhost:3000/oauth2/authorize?client_id=111-222-3333&response_type=token&redirect_uri=http://localhost:2000/oauth2_callback
</pre>

Note that you need to use the APIKey as the client_id parameter, and the redirect_uri registered for the application.

The url will present you with an authorization page. Fill in some credentials and click approve, or go straight to deny.
The browser will be redirected to your registered callback - if you use the supplied callback server it will display your
access token or the access denied error. You might also get an error if the client if can't be found, it it's not
a valid client id, etc.




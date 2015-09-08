## Roll - Services to maintain a roll of developers and their applications

Patterns and some code inspired by Hashicorp vault.

### Runtime set up

#### Vault

This software uses [Hashicorp's Vault](https://vaultproject.io/) for storing secret, in this case the private keys for signing tokens. We keep
the public keys in there are well, since we don't hand those out to anybody.

The easiest way to get set up with vault in a configuration that requires sealing/unsealing, use of tokens, and
so on is to run vault with the file backend.

To do use a configuration that looks like this:

<pre>
backend "file" {
        path = "/cygdrive/c/vault/backend"
}

listener "tcp" {
  address = "127.0.0.1:8200"
  tls_disable = 1
}
</pre>

Then run vault pointing to that configuration:

<pre>
vault server -config vconfig
</pre>

You will need to initialize the vault (`vault init`) and unseal the vault (`vault unseal`) using the key shards 
produced by the init process. Note - hang onto the key shards and root token! 

To interact with vault from the command line, you will need to set the VAULT_ADDR environment variable to
`http://localhost:8200` and set the VAULT_TOKEN environment variable to the root token (or make your owbn token).

Refer to the vault documentation for details.

#### AWS
You will need to set up AWS credentials for a user associated with the following policy:

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

The code assumes two tables have been created in DynamoDB:

* Developer with a string hash key named EMail
* Application with a string hash key named APIKey

A script to create these tables will be added to the project code eventually.

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
 
The following commands downloads all the dependencies needed to build the software:
 
<pre>
go get github.com/stretchr/testify/
go get -u github.com/aws/aws-sdk-go/...
go get github.com/hashicorp/vault/api
go get github.com/dgrijalva/jwt-go
go get github.com/nu7hatch/gouuid
</pre>


### Trying out the Implicit Grant Flow

#### Server Setup

Build the rollsvcs executable (go build in the rollsvcs directory), set the VAULT_ADDR and VAULT_TOKEN 
environment variables as described above, and start the server. Here we assume the use of port 3000.

<pre>
./rollsvcs -port 3000
</pre>

Next, build the callback server, which is used for the oauth2 callback for our sample, plus it can mock an
XTRAC login if you used xtrac://localhost:2000 as the login provider. We assume the use of port 2000.

<pre>
./cbserver -port 2000
</pre>

#### Data Setup

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
"RedirectURI":"http://localhost:2000/oauth2_callback",
"LoginProvider":"xtrac://localhost:2000"
}' localhost:3000/v1/applications/111-222-3333
</pre>

Retrieve registered app:

<pre>
curl -v localhost:3000/v1/applications/111-222-3333
</pre>


#### Executing the flow

Open the following in your browser. Note I use chrome for this - I assume it will work in other browsers.

<pre>
http://localhost:3000/oauth2/authorize?client_id=111-222-3333&response_type=token&redirect_uri=http://localhost:2000/oauth2_callback
</pre>

Note that you need to use the APIKey as the client_id parameter, and the redirect_uri registered for the application.

The url will present you with an authorization page. Fill in some credentials and click approve, or go straight to deny.
The browser will be redirected to your registered callback - if you use the supplied callback server it will display your
access token or the access denied error. You might also get an error if the client if can't be found, it it's not
a valid client id, etc.




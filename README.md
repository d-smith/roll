## Roll - Services to maintain a roll of developers and their applications

Patterns and some code inspired by Hashicorp vault.

<pre>
Currently in PoC status - don't even think about using this is a production situation.
</pre>

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

Also, if you need a proxy setting to access the internet, set the HTTP_PROXY environment variable.

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
"ApplicationName":"App No. 4",
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

Note in ths implementation of the flow, the credentials are sent back to the authorization server, which looks up
the authentication (login) endpoint associated with the application, which is embedded in a URL, e.g.

<pre>
xtrac://localhost:2000
</pre>

Using the URL scheme, we can accomodate different login providers, currently we know how to authenticate
against XTRAC.

### Authorization Code Flow

This can be be done with the above setup by modifying the above URL to use `code`
as the grant type:

<pre>
http://localhost:3000/oauth2/authorize?client_id=111-222-3333&response_type=code&redirect_uri=http://localhost:2000/oauth2_callback
</pre>

### Username Password Flow

This can be executed directly via curl, e.g.

<pre>
curl --data "client_id=111-222-3333" --data "grant_type=password" --data-urlencode "client_secret=EVNIFUt3hMFYb9aHy1N8LyEmTsLS3y+XK6xDvVbU+E0=" --data "username=foo" --data "password=passw0rd" localhost:3000/oauth2/token
</pre>

### JWT Flow

The JWT flow allows a security token created in a different fiefdom to be exchanged for an XTRAC token. To enable
this flow for an application, a certificate that can be used to validate the foreign JWT is uploaded to the roll server 
for the application. When the external token is posted to the token endpoint, the application key (client_id)
associated with the application is assumed to be carried in the token's iss claim: the public key extracted from the
uploaded certificate is used to validate the token signature, and if it checks out a access token is returned.

#### Trying it out

To try it out, we need to upload a cert, create and sign a JWT with the private key assocaited with the cert, and
post the JWT flow payload to the token endpoint -- jwt-sample/jwtflow-sample.go shows how to execute this flow.


### Protected Resource

Now that an application has been configured and an access token created, we can protect resources via
a simple wrapper.

The authzwrapper package contains a simple wrapper that restricts access to requests accompanied by authorization
bearer tokens created via the OAuth 2 flows supported by roll. 

The echo server provides an example of a protected resource.

To try it out, build the echo server and run it on say port 5000.

If you try it without a token, access will be denied:

<pre>
curl -X PUT -d 'Hello hello echo echo' localhost:5000/echo
Unauthorized
</pre>

If you use the token obtained through the implicit grant flow, access will be granted.

<pre>
curl -X PUT -d 'Hello hello echo echo' -H 'Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhcHBsaWNhdGlvbiI6IkhhbmdyeSBCaXJkeiIsImF1ZCI6IjExMS0yMjItMzMzMyIsImlhdCI6MTQ0MTcyNjA2NywianRpIjoiOGZjM2MzNzYtYjMyMy00Yzg5LTdiZTktOWRkZWE3ZWJhNWM2In0.yOjkodiiJtNnXGSoz2lipBgYNyKmQApjKVHPmkiW-peAVhtyQw-q3nnD-H93-vioiq-qvwKp9R4uj1gkPSXJlPJDDj4A6AtqlbbYElQ3K2q9IPPeYiaOR2fJZtLYsIvoDZimGHq_FjZvxDzYZalFSd7BDFeQ5xmhGWczqs6vNNE' localhost:5000/echo
Hello hello echo echo
</pre>

### TODO

* Create a token validation endpoint (to avoid the confused deputy problem)
* Keep track of the callback codes we generate to avoid replays while waiting for them to expire
* Make sure a callback code token can't be used as an access token
* Split the validate behavior out of the authz handler file
* Check auth code expiration
* What about refresh tokens?



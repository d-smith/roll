## Roll - Services to maintain a roll of developers and their applications

Patterns and some code inspired by Hashicorp vault.

<pre>
Currently in PoC status - don't even think about using this in a production situation.
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
  address = "0.0.0.0:8200"
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

#### AWS DynamoDB Support
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

The code uses two tables in DynamoDB: Application and Developer. Refer to the go
code under the repos/ddl package that can be used to create the tables.

Note that you can use DynamoDB Local as well with this code. To do set, specify local use via
an environment variable named LOCAL_DYNAMO_ADDR, and set your local address using this 
variable, e.g.

<pre>
export LOCAL_DYNAMO_ADDR=http://localhost:8000
</pre>

#### MariaDB Support

As an alternative to DynamoDB, Roll can be configured to use MariaDB as its database. To set up MariaDB,
use the rolldb.sql script in repos/ddl to create the database and user for the application. Minimally, the user
should be customized. Once the database and user have been created, the tabledefs.sql script can be run to create
the tables and grant the appropriate accesss to the roll db user created in the first step, which means tabledefs.sql
will need to be edited to reference the user created in step 1.

### Build Dependencies

Still need to vendor my dependencies, but they are:

* [Stretchr Testify](https://github.com/stretchr/testify/)
* [AWS Golang SDK](https://github.com/aws/aws-sdk-go)
* [Hashicorp Vault Golang API](https://github.com/hashicorp/vault/tree/master/api)
* [JWT-go](https://github.com/dgrijalva/jwt-go)
* [Go UUID](https://github.com/nu7hatch/gouuid)
* Context package from the [Gorilla Web Toolkit]()
* [Go MySQL Driver](github.com/go-sql-driver/mysql)

Use `go get github.com/hashicorp/vault/api` to install the API portion of Vault

I used the [mockery tool](https://github.com/vektra/mockery) to generate the mocks - I don't think there's a runtime
 dependency but go get it as per its instructions if you see build weirdness.
 
The following commands downloads all the dependencies needed to build the software:
 
<pre>
go get github.com/stretchr/testify/
go get -u github.com/aws/aws-sdk-go/...
go get github.com/hashicorp/vault/api
go get github.com/dgrijalva/jwt-go
go get github.com/nu7hatch/gouuid
go get github.com/gorilla/context
go get github.com/go-sql-driver/mysql
</pre>

### Test Dependencies

To run the integration test, install gucumber:

<pre>
go get github.com/lsegal/gucumber/cmd/gucumber
</pre>

You will also need the docker client to run the integration tests.

<pre>
go get github.com/samalba/dockerclient
</pre>

Running the integration tests requires docker-machine to be installed. You also need AWS credentials and, if
behind an http_proxy, set your http_proxy and no_proxy environment variables appropriately.

### Trying things out

Below are some older instructions (this documentation is in the process of being rewritten).

In a nutshell, you need to:

1. Boot roll in unsecure mode
    * In rollsvcs/cmd set the environment using setenv.sh then run `go run rollmain.go -port 3000 -unsecure`
2. Seed a roll application - see seed.js in the [Roll Setup](https://github.com/xtraclabs/rollsetup) repository. Seed.js
registers a developer and seeds a portal application associated with that developer. Roll can then being run in secure
mode as the registered application created by seed.js
3. Reboot roll in secure mode, using the application client_id obtained in the previous step.
3. Register a sample application for use with rollsample and rollecho in trying out the grants as outlined below. See
register.js in the [Roll Setup](https://github.com/xtraclabs/rollsetup) repository. Note you'll need to grab the
client id and secret for the seeded roll app from the application table and set the values of 
portClientId and clientSecret in register.js
4. Start the [rollsample](https://github.com/xtraclabs/rollsample) and 
[rollecho](https://github.com/xtraclabs/rollecho) applications. Roll sample is a  sample application that obtains auth tokens
 via roll, and roll echo is a service that requires an auth token associated with the sample application for
 admission.
5. Do the flows as explained below.

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

#### API Examples

The following command line examples using curl illustrate how to secure access to roll services to a specific application
registered with roll.

It's a two step process. First roll is booted in unsecured mode to allow an initial 'developer' to be created, followed by
the creation of the application that will be authorized to call roll. After this has been done, roll is booted in secure
mode, afterwhich all access will require an access token obtained using the authorized application.

Note the OAuth endpoints are not restricted, but obtaining access tokens is done in the context of an authorized
application, which is the source of the client id, client secret, etc.

#####Unsecured (Bootstrap)

<pre>
go run rollmain.go -port 3000 -unsecure
</pre>

<pre>
curl -v -X PUT -d '
{
"email":"test@dev.com",
"firstName":"Doug",
"lastName":"Dev"
}' -H 'X-Roll-Subject: portal-admin' localhost:3000/v1/developers/test@dev.com
</pre>

<pre>
curl -X POST -d '{
"applicationName":"secured dev portal",
"developerEmail":"test@dev.com",
"redirectURI":"http://localhost:2000/oauth2_callback",
"loginProvider":"xtrac://localhost:2000"
}' -H 'X-Roll-Subject: portal-admin' localhost:3000/v1/applications
{"client_id":"14fc366b-9138-46af-4bf7-c47579911028"}

curl  -H 'X-Roll-Subject: portal-admin' localhost:3000/v1/applications/14fc366b-9138-46af-4bf7-c47579911028
{"developerEmail":"test@dev.com","DeveloperID":"portal-admin","clientID":"14fc366b-9138-46af-4bf7-c47579911028","applicationName":"secured dev portal","clientSecret":"wW2vFB3aWjOkcj9LoONIpmT+VQ3/KoS8GiRom0yqReo=","redirectURI":"http://localhost:2000/oauth2_callback","loginProvider":"xtrac://localhost:2000","jwtFlowPublicKey":""}

curl --data "client_id=14fc366b-9138-46af-4bf7-c47579911028" --data "grant_type=password" --data-urlencode "client_secret=wW2vFB3aWjOkcj9LoONIpmT+VQ3/KoS8GiRom0yqReo=" --data "username=foo" --data "password=passw0rd" localhost:3000/oauth2/token
{"access_token":"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhcHBsaWNhdGlvbiI6InNlY3VyZWQgZGV2IHBvcnRhbCIsImF1ZCI6IjE0ZmMzNjZiLTkxMzgtNDZhZi00YmY3LWM0NzU3OTkxMTAyOCIsImV4cCI6MTQ1MjI4ODc3NCwiaWF0IjoxNDUyMjAyMzc0LCJqdGkiOiJkNWQ1ZWNhYy05ZWU3LTQ0YzctNmZjYS1kMDBhN2RjZTQ3ZmYiLCJzY29wZSI6IiIsInN1YiI6ImZvbyJ9.mgPyalqccJpA_IC2-JB3WcdqpRbxb2ZDkdNAeTqaNnMgxCbw1z6BIjOdpmHiZYc51g5zbTWUJl6XZsyv1Ul3u0U8kpfhjmg12VqN1mxVbsx3Z6wVXKRAaFN85rnoE-h4TaXRJcqq5wXzBXz1QT0qX8462VJN4cYaHo-MmHydTtU","token_type":"Bearer"}
</pre>


#####Secured

Note - use admins.go in repos/util to seed admin users - required to use a scope of admin.

<pre>
export ROLL_CLIENTID=1d703e17-fc84-42eb-65b6-9dcb7700b282
go run rollmain.go -port 3000
</pre>

<pre>
curl -v -X PUT -d '
{
"email":"new-dev@dev.com",
"firstName":"Doug",
"lastName":"Dev"
}' -H 'X-Roll-Subject: foo' localhost:3000/v1/developers/doug@dev.com

< HTTP/1.1 401 Unauthorized


curl --data "client_id=14fc366b-9138-46af-4bf7-c47579911028" --data "grant_type=password" --data-urlencode "client_secret=wW2vFB3aWjOkcj9LoONIpmT+VQ3/KoS8GiRom0yqReo=" --data "username=newdev" --data "password=passw0rd" localhost:3000/oauth2/token
{"access_token":"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhcHBsaWNhdGlvbiI6InNlY3VyZWQgZGV2IHBvcnRhbCIsImF1ZCI6IjY3YmZmYzQwLWVlZDItNDE0ZC01OWUzLWNhZTgyYzg3YzUyMCIsImV4cCI6MTQ1MjI4Njk4NiwiaWF0IjoxNDUyMjAwNTg2LCJqdGkiOiJjYWI2ZjE3OC0xNzA3LTQ2NDMtNGU3Mi03OWUzMDM1NDQwY2QiLCJzY29wZSI6IiIsInN1YiI6Im5ld2RldiJ9.lDWPsru75l-KdchoTZQE63LWaNX6f2dRwXlgWqLXXvbcd76R05zOKutOGUDQadtTTTudoRP4KQ2YHqxttPtVlnk9VSajJEuCysTdfslyegoTsc7qcy3vjxIc17FH6mclP8aqcgjyR20_L33GC7HVKe5CPRjOr615afVZF9GvXNs","token_type":"Bearer"}

curl -v -X PUT -d '
{
"email":"new-dev@dev.com",
"firstName":"Doug",
"lastName":"Dev"
}' -H 'Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhcHBsaWNhdGlvbiI6InNlY3VyZWQgZGV2IHBvcnRhbCIsImF1ZCI6IjE0ZmMzNjZiLTkxMzgtNDZhZi00YmY3LWM0NzU3OTkxMTAyOCIsImV4cCI6MTQ1MjI4ODkwMSwiaWF0IjoxNDUyMjAyNTAxLCJqdGkiOiJiNjdmYTFiOS1kM2VlLTQyYTMtNDIzNS03Yzg4YmMyNTFmYWYiLCJzY29wZSI6IiIsInN1YiI6Im5ld2RldiJ9.EReJee2p__3KCvzFIj7esj5rGZy0BSWzQ24gwDDA4yFehCcvfPlWMP3M4_31tQeNJjgzo4PxfyfIvHTfKvbZy3h4OJru7Rk9ECkgMHx3yM-mWVKGvGJ3xdHnPcAbT8ArDcQRthS_So5KYJ5I3hq_swc7rOH0wdoF0FDMjD2fuP8' localhost:3000/v1/developers/new-dev@dev.com


curl -H 'Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhcHBsaWNhdGlvbiI6InNlY3VyZWQgZGV2IHBvcnRhbCIsImF1ZCI6IjE0ZmMzNjZiLTkxMzgtNDZhZi00YmY3LWM0NzU3OTkxMTAyOCIsImV4cCI6MTQ1MjI4ODkwMSwiaWF0IjoxNDUyMjAyNTAxLCJqdGkiOiJiNjdmYTFiOS1kM2VlLTQyYTMtNDIzNS03Yzg4YmMyNTFmYWYiLCJzY29wZSI6IiIsInN1YiI6Im5ld2RldiJ9.EReJee2p__3KCvzFIj7esj5rGZy0BSWzQ24gwDDA4yFehCcvfPlWMP3M4_31tQeNJjgzo4PxfyfIvHTfKvbZy3h4OJru7Rk9ECkgMHx3yM-mWVKGvGJ3xdHnPcAbT8ArDcQRthS_So5KYJ5I3hq_swc7rOH0wdoF0FDMjD2fuP8' localhost:3000/v1/developers/new-dev@dev.com

</pre>

<pre>
curl -X POST -d '{
"applicationName":"App No. 5",
"developerEmail":"new-dev@dev.com",
"redirectURI":"http://localhost:2000/oauth2_callback",
"loginProvider":"xtrac://localhost:2000"
}' -H 'Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhcHBsaWNhdGlvbiI6InNlY3VyZWQgZGV2IHBvcnRhbCIsImF1ZCI6IjE0ZmMzNjZiLTkxMzgtNDZhZi00YmY3LWM0NzU3OTkxMTAyOCIsImV4cCI6MTQ1MjI4ODkwMSwiaWF0IjoxNDUyMjAyNTAxLCJqdGkiOiJiNjdmYTFiOS1kM2VlLTQyYTMtNDIzNS03Yzg4YmMyNTFmYWYiLCJzY29wZSI6IiIsInN1YiI6Im5ld2RldiJ9.EReJee2p__3KCvzFIj7esj5rGZy0BSWzQ24gwDDA4yFehCcvfPlWMP3M4_31tQeNJjgzo4PxfyfIvHTfKvbZy3h4OJru7Rk9ECkgMHx3yM-mWVKGvGJ3xdHnPcAbT8ArDcQRthS_So5KYJ5I3hq_swc7rOH0wdoF0FDMjD2fuP8' localhost:3000/v1/applications
{"client_id":"aaf9404c-ac51-4720-6097-f80c5404773d"}

curl -H 'Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhcHBsaWNhdGlvbiI6InNlY3VyZWQgZGV2IHBvcnRhbCIsImF1ZCI6IjE0ZmMzNjZiLTkxMzgtNDZhZi00YmY3LWM0NzU3OTkxMTAyOCIsImV4cCI6MTQ1MjI4ODkwMSwiaWF0IjoxNDUyMjAyNTAxLCJqdGkiOiJiNjdmYTFiOS1kM2VlLTQyYTMtNDIzNS03Yzg4YmMyNTFmYWYiLCJzY29wZSI6IiIsInN1YiI6Im5ld2RldiJ9.EReJee2p__3KCvzFIj7esj5rGZy0BSWzQ24gwDDA4yFehCcvfPlWMP3M4_31tQeNJjgzo4PxfyfIvHTfKvbZy3h4OJru7Rk9ECkgMHx3yM-mWVKGvGJ3xdHnPcAbT8ArDcQRthS_So5KYJ5I3hq_swc7rOH0wdoF0FDMjD2fuP8' localhost:3000/v1/applications/aaf9404c-ac51-4720-6097-f80c5404773d
{"developerEmail":"new-dev@dev.com","DeveloperID":"newdev","clientID":"aaf9404c-ac51-4720-6097-f80c5404773d","applicationName":"App No. 5","clientSecret":"JqBUEAYXOxook/VzwqiqyiKPirJ5mBXZkOerrhAfId8=","redirectURI":"http://localhost:2000/oauth2_callback","loginProvider":"xtrac://localhost:2000","jwtFlowPublicKey":""}
</pre>

<pre>
curl --data "client_id=aaf9404c-ac51-4720-6097-f80c5404773d" --data scope=admin --data "grant_type=password" --data-urlencode "client_secret=JqBUEAYXOxook/VzwqiqyiKPirJ5mBXZkOerrhAfId8=" --data "username=portal-admin" --data "password=passw0rd" localhost:3000/oauth2/token
{"access_token":"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhcHBsaWNhdGlvbiI6IkFwcCBOby4gNSIsImF1ZCI6ImFhZjk0MDRjLWFjNTEtNDcyMC02MDk3LWY4MGM1NDA0NzczZCIsImV4cCI6MTQ1MjI4OTUwOSwiaWF0IjoxNDUyMjAzMTA5LCJqdGkiOiIzODMxMDVhNy1kODMwLTRkZjgtNzQwNy1iMmMyZjZiYzFkMjAiLCJzY29wZSI6ImFkbWluIiwic3ViIjoicG9ydGFsLWFkbWluIn0.k8LgLK1OPDCeN8XzByS80UePhhB16CdnXPBIGdG6lRbJ97qyPMticFNuzvPCQq2ZKJ4Qp1gvIDYEFcq6iawWcMLdl_tuyhGT-wkq3_yeZq2bk2aKmbkws4azKhNOv8Ih-bYg4E7PjVuowplkyUvUqm7_tBqKElFQW-3kYBYkyyk","token_type":"Bearer"}

curl -H 'Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhcHBsaWNhdGlvbiI6InNlY3VyZWQgZGV2IHBvcnRhbCIsImF1ZCI6IjE0ZmMzNjZiLTkxMzgtNDZhZi00YmY3LWM0NzU3OTkxMTAyOCIsImV4cCI6MTQ1MjI4ODc3NCwiaWF0IjoxNDUyMjAyMzc0LCJqdGkiOiJkNWQ1ZWNhYy05ZWU3LTQ0YzctNmZjYS1kMDBhN2RjZTQ3ZmYiLCJzY29wZSI6IiIsInN1YiI6ImZvbyJ9.mgPyalqccJpA_IC2-JB3WcdqpRbxb2ZDkdNAeTqaNnMgxCbw1z6BIjOdpmHiZYc51g5zbTWUJl6XZsyv1Ul3u0U8kpfhjmg12VqN1mxVbsx3Z6wVXKRAaFN85rnoE-h4TaXRJcqq5wXzBXz1QT0qX8462VJN4cYaHo-MmHydTtU' localhost:3000/v1/applications
</pre>

##### Previous Example
Create Developer

<pre>
curl -v -X PUT -d '
{
"email":"doug@dev.com",
"firstName":"Doug",
"lastName":"Dev"
}' localhost:3000/v1/developers/doug@dev.com
</pre>

Retrieve a Developer

<pre>
curl localhost:3000/v1/developers/doug@dev.com
{"firstName":"Doug","lastName":"Dev","email":"doug@dev.com","id":""}
</pre>

List Developers

<pre>
curl localhost:3000/v1/developers
[{"firstName":"Doug","lastName":"Dev","email":"doug@dev.com","id":""}]
</pre>

Register an application

<pre>
curl -X POST -d '{
"applicationName":"App No. 5",
"developerEmail":"doug@dev.com",
"redirectURI":"http://localhost:2000/oauth2_callback",
"loginProvider":"xtrac://localhost:2000"
}' localhost:3000/v1/applications
{"client_id":"7843541e-d4cb-4903-5b88-ee596c32ecd7"}
</pre>

Update an application

<pre>
curl -v -X PUT -d '{
"applicationName":"App No. Four",
"developerEmail":"doug@dev.com",
"redirectURI":"http://localhost:2000/oauth2_callback",
"loginProvider":"xtrac://localhost:2000"
}' localhost:3000/v1/applications/7843541e-d4cb-4903-5b88-ee596c32ecd7
</pre>

Retrieve an Application

<pre>
curl localhost:3000/v1/applications/7843541e-d4cb-4903-5b88-ee596c32ecd7
{"developerEmail":"doug@dev.com","clientID":"7843541e-d4cb-4903-5b88-ee596c32ecd7","applicationName":"App No. Four","clientSecret":"bQeH+n/Q9g8gM++Xd9gnqrn6zp92EZpSXrRPofVUbyk=","redirectURI":"http://localhost:2000/oauth2_callback","loginProvider":"xtrac://localhost:2000","jwtFlowPublicKey":""}
</pre>

Retrieve all applications

<pre>
curl localhost:3000/v1/applications
[{"developerEmail":"doug@dev.com","clientID":"7843541e-d4cb-4903-5b88-ee596c32ecd7","applicationName":"App No. Four","clientSecret":"bQeH+n/Q9g8gM++Xd9gnqrn6zp92EZpSXrRPofVUbyk=","redirectURI":"http://localhost:2000/oauth2_callback","loginProvider":"xtrac://localhost:2000","jwtFlowPublicKey":""}]
</pre>

#### Executing the flow

Open the following in your browser. Note I use chrome for this - I assume it will work in other browsers.

<pre>
http://localhost:3000/oauth2/authorize?client_id=111-222-3333&response_type=token&redirect_uri=http://localhost:2000/oauth2_callback
</pre>

Note that you need to use the ClientID as the client_id parameter, and the redirect_uri registered for the application.

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
curl --data "client_id=7843541e-d4cb-4903-5b88-ee596c32ecd7" --data "grant_type=password" --data-urlencode "client_secret=bQeH+n/Q9g8gM++Xd9gnqrn6zp92EZpSXrRPofVUbyk=" --data "username=foo" --data "password=passw0rd" localhost:3000/oauth2/token
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



To run the sample, first get an auth token via username/password flow. 

Use the username of the developer who owns the application.

curl --data "client_id=d80c9758-f965-4286-7003-e5af065d8082" --data "grant_type=password" --data-urlencode "client_secret=4suKio+75Vsm2o+i8H8dEjxVbXhoB81G+rJolRw8+XU=" --data "username=user" --data "password=passw0rd" localhost:3000/oauth2/token

Paste the returned token into the jwt-sample code as the value of authToken

Set client id and client secret for the app that is being updated.
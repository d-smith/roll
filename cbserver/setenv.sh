#!/bin/sh

# This script sets up the application parameters for the auth code
# callback sample. The sample assumes it's an application that
# was registered with the application roll.

export AZ_SERVER=localhost:3000
export CLIENT_ID=78e1cca4-d228-464d-4b0c-492311bf3b73
export CLIENT_SECRET=Pw/QVJWAzOViuV8Z4tttryjnSsSJ4QI8Mi5Tg+LnetM=
export REDIRECT_URI=http://localhost:2000/oauth2_callback
export ECHO_ENDPOINT=http://localhost:5000

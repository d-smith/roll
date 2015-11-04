#!/bin/sh

# This script sets up the application parameters for the auth code
# callback sample. The sample assumes it's an application that
# was registered with the application roll.

export AZ_SERVER=localhost:3000
export CLIENT_ID=3ca926b9-44eb-4ef2-7971-aa33b1620f78
export CLIENT_SECRET=VoscDd4uj22UhXrwe++RceNeDqJZ0ZZwN8PBMS1BUlM=
export REDIRECT_URI=http://localhost:2000/oauth2_callback
export ECHO_ENDPOINT=http://localhost:5000

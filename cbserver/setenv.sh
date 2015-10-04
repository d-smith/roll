#!/bin/sh

# This script sets up the application parameters for the auth code
# callback sample. The sample assumes it's an application that
# was registered with the application roll.

export AZ_SERVER=localhost:3000
export CLIENT_ID=8a112572-47bb-4c95-5179-3c26fa0acd27
export CLIENT_SECRET=8fnyNejHEMiAn5jXlRVuShq+AaAaEzU/HrfTvW61Tvw=
export REDIRECT_URI=http://localhost:2000/oauth2_callback
export ECHO_ENDPOINT=http://localhost:5000

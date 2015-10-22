#!/bin/sh

# This script sets up the application parameters for the auth code
# callback sample. The sample assumes it's an application that
# was registered with the application roll.

export AZ_SERVER=localhost:3000
export CLIENT_ID=f12c5788-14b6-4082-41c8-9c8b358feb65
export CLIENT_SECRET=aIstHPhjRsvc2ZgkR34xinh4BH+qlt4h4BDw/1tgxCE=
export REDIRECT_URI=http://localhost:2000/oauth2_callback
export ECHO_ENDPOINT=http://localhost:5000

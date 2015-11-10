#!/bin/sh

# This script sets up the application parameters for the auth code
# callback sample. The sample assumes it's an application that
# was registered with the application roll.

export AZ_SERVER=MACLB015803:3000
export CLIENT_ID=23ce181a-7177-4ac5-554e-11151a1ff342
export CLIENT_SECRET=fhoD5vM2suoaqpgF0+L4yjR4hg8sRnw/GsIejXutvmc=
export REDIRECT_URI=http://MACLB015803:2000/oauth2_callback
export ECHO_ENDPOINT=http://MACLB015803:5000

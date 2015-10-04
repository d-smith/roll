#!/bin/sh

# This script sets up the application parameters for the auth code
# callback sample. The sample assumes it's an application that
# was registered with the application roll.

export AZ_SERVER=localhost:3000
export CLIENT_ID=399fe4c3-c4d0-4a14-7c4e-fea81e060229
export CLIENT_SECRET=a3BAA8SZ9hXEKi+U9Ut9s1xEdsH8OvZAWx0ATnwDR7c=
export REDIRECT_URI=http://localhost:2000/oauth2_callback

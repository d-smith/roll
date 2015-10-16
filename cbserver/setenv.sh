#!/bin/sh

# This script sets up the application parameters for the auth code
# callback sample. The sample assumes it's an application that
# was registered with the application roll.

export AZ_SERVER=localhost:3000
export CLIENT_ID=86729cb4-9e40-4277-6755-8c0192a97306
export CLIENT_SECRET=5OrhMgtvN5kiHyRuA61WkY5S22QWOXxWx7yWQR14C/s=
export REDIRECT_URI=http://localhost:2000/oauth2_callback
export ECHO_ENDPOINT=http://localhost:5000

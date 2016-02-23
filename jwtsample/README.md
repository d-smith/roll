##Token Exchange

The basic set up is the signing cert for external JWTs to be exchanged for roll JWTs is registered for
an application, along with the iss identifying the token issuer, and the aud identifying the intended
recipient of JWT.

When receiving a token to exchange via the OAuth2 JWT flow, Roll looks up the application definition 
associated with the audience (which implies each registered application must be accessed with an intended
audience for external access). The signing cert stored for the application is used to validate the token
signature, and we also check the issuer of the token matches the issuer we store with the application
definition.

If the foreign token checks out, we issue a roll JWT with the sub claim conveyed by the foreign token.

Note that currently we only support the RS256 signing method, to avoid the potential vulnerabilities with
HS256.

To test this, we used set up an application in auth0. We stored the auth0 client id, the signing cert, and the 
issuing domain in Roll.

When presented a auth0 token containing the following:

<pre>
{
  "iss": "https://xavi.auth0.com/",
  "sub": "auth0|56cc9aae7541fce12c63247c",
  "aud": "vY0bFoxCBzE9rrTNTEjhIfay8MbFYq9Z",
  "exp": 1456285522,
  "iat": 1456249522
}
</pre>

We exchanged it for a token containing this:

<pre>
{
  "application": "dev portal",
  "aud": "5d130f17-2fe5-4462-4e9d-9b6eb2d806e8",
  "exp": 1456337947,
  "iat": 1456251547,
  "jti": "15e48228-f900-4b0f-7aa8-146b0f69edf9",
  "scope": "",
  "sub": "auth0|56cc9aae7541fce12c63247c"
}
</pre>
uaa-demo
========

## Purpose
* Playing with golang
* How to retrieve the token from cloudfoundry uaa through "auth code" flow
* Better understanding the uaa flow

## Explain

* Currently this demo shows how to use authcode authentication to retrieve token
* Using https://github.com/pivotal-cf/uaa-sso-golang/ for uaa client
* Using gorrila session for session management

## Flow (Keep it as simple as possible)

* /view shows the token and refresh value

  If the token is not in the session, redirect to the login server of cloudfoundry.
  After user input username and password with login server, it redirects to the "/token" handler with the code
* /token

  This handler retrieve the token through the code. set the token to the session and redirect back to /view

## Note: 
* Configure config.json for the uaa information
* "cf push" with manifest file to push to cloudfoundry environment. Sample config.json is for a uaa environment


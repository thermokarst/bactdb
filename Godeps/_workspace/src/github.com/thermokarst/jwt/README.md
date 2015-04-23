# jwt

[![GoDoc](https://godoc.org/github.com/thermokarst/jwt?status.svg)](https://godoc.org/github.com/thermokarst/jwt)

A simple, opinionated Go net/http middleware for integrating JSON Web Tokens into
your application:

```go
package main

import (
    "errors"
    "fmt"
    "net/http"
    "time"

    "github.com/thermokarst/jwt"
)

func protectMe(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "secured")
}

func main() {
    var authFunc = func(email string, password string) error {
        // Hard-code a user --- this could easily be a database call, etc.
        if email != "test" || password != "test" {
            return errors.New("invalid credentials")
        }
        return nil
    }

    var claimsFunc = func(userId string) (map[string]interface{}, error) {
        currentTime := time.Now()
        return map[string]interface{}{
            "iat": currentTime.Unix(),
            "exp": currentTime.Add(time.Minute * 60 * 24).Unix(),
            "sub": userId,
        }, nil
    }

    var verifyClaimsFunc = func(claims []byte) error {
        currentTime := time.Now()
        var c struct {
            Exp int64
            Iat int64
            Sub string
        }
        err := json.Unmarshal(claims, &c)
        if err != nil {
            return err
        }
        if currentTime.After(time.Unix(c.Exp, 0)) {
            return errors.New("this token has expired!")
        }
        if c.Sub != "test" {
            return errors.New("who are you??!")
        }
        return nil
    }

    config := &jwt.Config{
        Secret: "password",
        Auth:   authFunc,
        Claims: claimsFunc,
    }
    j, err := jwt.New(config)
    if err != nil {
        panic(err)
    }
    protect := http.HandlerFunc(protectMe)
    http.Handle("/authenticate", j.GenerateToken())
    http.Handle("/secure", j.Secure(protect, verifyClaimsFunc))
    http.ListenAndServe(":8080", nil)
}
```

# Installation

    $ go get github.com/thermokarst/jwt

# Usage

**This is a work in progress**

Create a new instance of the middleware by passing in a configuration for your
app.  The config includes a shared secret (this middleware only builds HS256
tokens), a function for authenticating user, and a function for generating a
user's claims. The idea here is to be dead-simple for someone to drop this into
a project and hit the ground running.

```go
config := &jwt.Config{
    Secret: "password",
    Auth:   authFunc, // func(string, string) error
    Claims: claimsFunc, // func(string) (map[string]interface{})
}
j, err := jwt.New(config)
```

Once the middleware is instantiated, create a route for users to generate a JWT
at.

```go
http.Handle("/authenticate", j.GenerateToken())
```

The auth function takes two arguments (the identity, and the authorization
key), POSTed as a JSON-encoded body:

    {"email":"user@example.com","password":"mypassword"}

These fields are static for now, but will be customizable in a later release.
The claims are generated using the claims function provided in the
configuration. This function is only run if the auth function verifies the
user's identity, then the user's unique identifier (primary key id, UUID,
email, whatever you want) is passed as a string to the claims function. Your
function should return a `map[string]interface{}` with the desired claimset.

Routes are "secured" by calling the `Secure(http.Handler, jwt.VerifyClaimsFunc)`
handler:

```go
http.Handle("/secureendpoint", j.Secure(someHandler, verifyClaimsFunc))
```

The claims verification function is called after the token has been parsed and
validated: this is where you control how your application handles the claims
contained within the JWT.

# Motivation

This work was prepared for a crypto/security class at the University of Alaska
Fairbanks.  I hope to use this in some of my projects, but please proceed with
caution if you adopt this for your own work. As well, the API is still quite
unstable, so be prepared for handling any changes.

# Tests

    $ go test

# Contributors

Matthew Ryan Dillon (matthewrdillon@gmail.com)


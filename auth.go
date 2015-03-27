package main

import (
	"errors"
	"os"
)

const (
	tokenName = "AccessToken"
)

var (
	verifyKey, signKey    []byte
	errWhileSigningToken  = errors.New("error while signing token")
	errPleaseLogIn        = errors.New("please log in")
	errWhileParsingCookie = errors.New("error while parsing cookie")
	errTokenExpired       = errors.New("token expired")
	errGenericError       = errors.New("generic error")
	errAccessDenied       = errors.New("insufficient privileges")
)

func setupCerts() error {
	// openssl genrsa -out app.rsa 1024
	signkey := os.Getenv("PRIVATE_KEY")
	if signkey == "" {
		return errors.New("please set PRIVATE_KEY")
	}
	signKey = []byte(signkey)

	// openssl rsa -in app.rsa -pubout > app.rsa.pub
	verifykey := os.Getenv("PUBLIC_KEY")
	if verifykey == "" {
		return errors.New("please set PUBLIC_KEY")
	}
	verifyKey = []byte(verifykey)

	return nil
}

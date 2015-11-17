package auth

import (
	"os"
	"time"

	"github.com/thermokarst/bactdb/Godeps/_workspace/src/github.com/thermokarst/jwt"
	"github.com/thermokarst/bactdb/models"
)

var (
	// Middleware is for JWT
	Middleware *jwt.Middleware
	// Config handles JWT middleware configuration
	Config = &jwt.Config{
		Secret:        os.Getenv("SECRET"),
		Auth:          models.DbAuthenticate,
		Claims:        claimsFunc,
		IdentityField: "username",
		VerifyField:   "password",
	}
)

func claimsFunc(email string) (map[string]interface{}, error) {
	// TODO: use helper
	currentTime := time.Now()
	user, err := models.DbGetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"name": user.Name,
		"iss":  "bactdb",
		"sub":  user.ID,
		"role": user.Role,
		"iat":  currentTime.Unix(),
		"exp":  currentTime.Add(time.Minute * 60 * 24).Unix(),
		"ref":  "",
	}, nil
}

func init() {
	var err error
	Middleware, err = jwt.New(Config)
	if err != nil {
		panic(err)
	}
}

package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/thermokarst/bactdb/Godeps/_workspace/src/github.com/gorilla/context"
	"github.com/thermokarst/bactdb/Godeps/_workspace/src/github.com/thermokarst/jwt"
	"github.com/thermokarst/bactdb/auth"
	"github.com/thermokarst/bactdb/errors"
	"github.com/thermokarst/bactdb/helpers"
	"github.com/thermokarst/bactdb/models"
	"github.com/thermokarst/bactdb/types"
)

func verifyClaims(claims []byte, r *http.Request) error {
	currentTime := time.Now()
	var c types.Claims
	err := json.Unmarshal(claims, &c)
	if err != nil {
		return err
	}

	if currentTime.After(time.Unix(c.Exp, 0)) {
		return errors.ErrExpiredToken
	}

	user, err := models.GetUser(c.Sub, "", &c)
	if err != nil {
		return err
	}

	if c.Role != user.Role {
		return errors.ErrInvalidToken
	}

	context.Set(r, "claims", c)
	return nil
}

func tokenHandler(h http.Handler) http.Handler {
	token := func(w http.ResponseWriter, r *http.Request) {
		recorder := httptest.NewRecorder()
		h.ServeHTTP(recorder, r)

		for key, val := range recorder.Header() {
			w.Header()[key] = val
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(recorder.Code)

		tokenData := string(recorder.Body.Bytes())

		var data []byte

		if recorder.Code != 200 {
			data, _ = json.Marshal(struct {
				Error string `json:"error"`
			}{
				Error: tokenData,
			})
		} else {
			data, _ = json.Marshal(struct {
				Token string `json:"access_token"`
			}{
				Token: tokenData,
			})
		}

		w.Write(data)
		return

	}
	return http.HandlerFunc(token)
}

func tokenRefresh(j *jwt.Middleware) errorHandler {
	t := func(w http.ResponseWriter, r *http.Request) *types.AppError {
		claims := helpers.GetClaims(r)
		user, err := models.GetUser(claims.Sub, "", &claims)
		if err != nil {
			return newJSONError(err, http.StatusInternalServerError)
		}

		user.Password = ""
		token, err := auth.Middleware.CreateToken(user.Email)
		if err != nil {
			return newJSONError(err, http.StatusInternalServerError)
		}

		data, _ := json.Marshal(struct {
			Token string `json:"access_token"`
		}{
			Token: token,
		})

		w.Write(data)
		return nil
	}
	return t
}

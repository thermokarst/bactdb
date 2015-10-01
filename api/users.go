package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/thermokarst/bactdb/Godeps/_workspace/src/github.com/gorilla/mux"
	"github.com/thermokarst/bactdb/Godeps/_workspace/src/github.com/lib/pq"
	"github.com/thermokarst/bactdb/Godeps/_workspace/src/github.com/mailgun/mailgun-go"
	"github.com/thermokarst/bactdb/Godeps/_workspace/src/golang.org/x/crypto/bcrypt"
	"github.com/thermokarst/bactdb/auth"
	"github.com/thermokarst/bactdb/helpers"
	"github.com/thermokarst/bactdb/models"
	"github.com/thermokarst/bactdb/payloads"
	"github.com/thermokarst/bactdb/types"
)

var (
	// TODO: fix this
	ErrUserNotFoundJSON      = types.NewJSONError(models.ErrUserNotFound, http.StatusNotFound)
	ErrUserNotUpdatedJSON    = types.NewJSONError(models.ErrUserNotUpdated, http.StatusBadRequest)
	ErrEmailAddressTakenJSON = types.NewJSONError(models.ErrEmailAddressTaken, http.StatusBadRequest)
	MgAccts                  = make(map[string]mailgun.Mailgun)
)

type UserService struct{}

func (u UserService) Unmarshal(b []byte) (types.Entity, error) {
	var uj payloads.UserPayload
	err := json.Unmarshal(b, &uj)
	return &uj, err
}

func (u UserService) List(val *url.Values, claims *types.Claims) (types.Entity, *types.AppError) {
	if val == nil {
		return nil, helpers.ErrMustProvideOptionsJSON
	}
	var opt helpers.ListOptions
	if err := helpers.SchemaDecoder.Decode(&opt, *val); err != nil {
		return nil, types.NewJSONError(err, http.StatusInternalServerError)
	}

	// TODO: fix this
	users := make(models.Users, 0)
	sql := `SELECT id, email, 'password' AS password, name, role,
		created_at, updated_at, deleted_at
		FROM users
		WHERE verified IS TRUE
		AND deleted_at IS NULL;`
	if err := models.DBH.Select(&users, sql); err != nil {
		return nil, types.NewJSONError(err, http.StatusInternalServerError)
	}
	return &users, nil
}

func (u UserService) Get(id int64, dummy string, claims *types.Claims) (types.Entity, *types.AppError) {
	user, err := models.DbGetUserById(id)
	user.Password = ""
	if err != nil {
		return nil, types.NewJSONError(err, http.StatusInternalServerError)
	}

	user.CanEdit = claims.Role == "A" || id == claims.Sub

	payload := payloads.UserPayload{
		User: user,
		Meta: &models.UserMeta{
			CanAdd: claims.Role == "A",
		},
	}
	return &payload, nil
}

func (u UserService) Update(id int64, e *types.Entity, dummy string, claims *types.Claims) *types.AppError {
	user := (*e).(*payloads.UserPayload).User

	original_user, err := models.DbGetUserById(id)
	if err != nil {
		return types.NewJSONError(err, http.StatusInternalServerError)
	}

	user.Id = id
	user.Password = original_user.Password
	user.Verified = original_user.Verified
	user.UpdatedAt = helpers.CurrentTime()

	if err := user.Validate(); err != nil {
		return &types.AppError{Error: err, Status: helpers.StatusUnprocessableEntity}
	}

	// TODO: fix this
	count, err := models.DBH.Update(user)
	user.Password = ""
	if err != nil {
		return types.NewJSONError(err, http.StatusInternalServerError)
	}
	if count != 1 {
		return ErrUserNotUpdatedJSON
	}

	return nil
}

func (u UserService) Create(e *types.Entity, dummy string, claims *types.Claims) *types.AppError {
	user := (*e).(*payloads.UserPayload).User
	if err := user.Validate(); err != nil {
		return &types.AppError{Error: err, Status: helpers.StatusUnprocessableEntity}
	}
	ct := helpers.CurrentTime()
	user.CreatedAt = ct
	user.UpdatedAt = ct
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return types.NewJSONError(err, http.StatusInternalServerError)
	}
	user.Password = string(hash)
	user.Role = "R"
	user.Verified = false

	// TODO: fix this
	if err := models.DBH.Insert(user); err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code == "23505" {
				return ErrEmailAddressTakenJSON
			}
		}
		return types.NewJSONError(err, http.StatusInternalServerError)
	}

	user.Password = "password" // don't want to send the hashed PW back to the client

	q := `INSERT INTO verification (user_id, nonce, referer, created_at) VALUES ($1, $2, $3, $4);`
	// TODO: move helpers.GenerateNonce
	nonce, err := helpers.GenerateNonce()
	if err != nil {
		return types.NewJSONError(err, http.StatusInternalServerError)
	}
	// TODO: fix this
	_, err = models.DBH.Exec(q, user.Id, nonce, claims.Ref, ct)
	if err != nil {
		return types.NewJSONError(err, http.StatusInternalServerError)
	}

	// Send out confirmation email
	// TODO: clean this up
	mg, ok := MgAccts[claims.Ref]
	if ok {
		sender := fmt.Sprintf("%s Admin <admin@%s>", mg.Domain(), mg.Domain())
		recipient := fmt.Sprintf("%s <%s>", user.Name, user.Email)
		subject := fmt.Sprintf("New Account Confirmation - %s", mg.Domain())
		message := fmt.Sprintf("You are receiving this message because this email "+
			"address was used to sign up for an account at %s. Please visit this "+
			"URL to complete the sign up process: %s/users/new/verify/%s. If you "+
			"did not request an account, please disregard this message.",
			mg.Domain(), claims.Ref, nonce)
		m := mailgun.NewMessage(sender, subject, message, recipient)
		_, _, err := mg.Send(m)
		if err != nil {
			log.Printf("%+v\n", err)
			return types.NewJSONError(err, http.StatusInternalServerError)
		}
	}

	return nil
}

func HandleUserVerify(w http.ResponseWriter, r *http.Request) *types.AppError {
	// TODO: clean this up
	nonce := mux.Vars(r)["Nonce"]
	q := `SELECT user_id, referer FROM verification WHERE nonce=$1;`

	var ver struct {
		User_id int64
		Referer string
	}
	if err := models.DBH.SelectOne(&ver, q, nonce); err != nil {
		log.Print(err)
		return types.NewJSONError(err, http.StatusInternalServerError)
	}

	if ver.User_id == 0 {
		return types.NewJSONError(errors.New("No user found"), http.StatusInternalServerError)
	}

	var user models.User
	if err := models.DBH.Get(&user, ver.User_id); err != nil {
		return types.NewJSONError(err, http.StatusInternalServerError)
	}

	user.UpdatedAt = helpers.CurrentTime()
	user.Verified = true

	count, err := models.DBH.Update(&user)
	if err != nil {
		return types.NewJSONError(err, http.StatusInternalServerError)
	}
	if count != 1 {
		return types.NewJSONError(errors.New("Count 0"), http.StatusInternalServerError)
	}

	q = `DELETE FROM verification WHERE user_id=$1;`
	_, err = models.DBH.Exec(q, user.Id)
	if err != nil {
		return types.NewJSONError(err, http.StatusInternalServerError)
	}
	fmt.Fprintln(w, `{"msg":"All set! Please log in."}`)
	return nil
}

func HandleUserLockout(w http.ResponseWriter, r *http.Request) *types.AppError {
	email := r.FormValue("email")
	if email == "" {
		return types.NewJSONError(errors.New("missing email"), http.StatusInternalServerError)
	}
	token, err := auth.Middleware.CreateToken(email)
	if err != nil {
		return types.NewJSONError(err, http.StatusInternalServerError)
	}
	origin := r.Header.Get("Origin")
	hostUrl, err := url.Parse(origin)
	if err != nil {
		return types.NewJSONError(err, http.StatusInternalServerError)
	}
	hostUrl.Path += "/users/lockoutauthenticate"
	params := url.Values{}
	params.Add("token", token)
	hostUrl.RawQuery = params.Encode()

	// Send out email
	// TODO: clean this up
	mg, ok := MgAccts[origin]
	if ok {
		sender := fmt.Sprintf("%s Admin <admin@%s>", mg.Domain(), mg.Domain())
		recipient := fmt.Sprintf("%s", email)
		subject := fmt.Sprintf("Password Reset Request - %s", mg.Domain())
		message := fmt.Sprintf("You are receiving this message because this email "+
			"address was used in an account lockout request at %s. Please visit "+
			"this URL to complete the process: %s. If you did not request help "+
			"with a lockout, please disregard this message.",
			mg.Domain(), hostUrl.String())
		m := mailgun.NewMessage(sender, subject, message, recipient)
		_, _, err := mg.Send(m)
		if err != nil {
			log.Printf("%+v\n", err)
			return types.NewJSONError(err, http.StatusInternalServerError)
		}
	}

	fmt.Fprintln(w, `{}`)
	return nil
}

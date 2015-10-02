package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/thermokarst/bactdb/Godeps/_workspace/src/github.com/gorilla/mux"
	"github.com/thermokarst/bactdb/Godeps/_workspace/src/github.com/lib/pq"
	"github.com/thermokarst/bactdb/Godeps/_workspace/src/github.com/mailgun/mailgun-go"
	"github.com/thermokarst/bactdb/Godeps/_workspace/src/golang.org/x/crypto/bcrypt"
	"github.com/thermokarst/bactdb/auth"
	"github.com/thermokarst/bactdb/errors"
	"github.com/thermokarst/bactdb/helpers"
	"github.com/thermokarst/bactdb/models"
	"github.com/thermokarst/bactdb/payloads"
	"github.com/thermokarst/bactdb/types"
)

var (
	// MgAccts is a map of Mailgun accounts.
	MgAccts = make(map[string]mailgun.Mailgun)
)

// UserService provides for CRUD operations.
type UserService struct{}

// Unmarshal satisfies interface Updater and interface Creater.
func (u UserService) Unmarshal(b []byte) (types.Entity, error) {
	var uj payloads.User
	err := json.Unmarshal(b, &uj)
	return &uj, err
}

// List lists all users.
func (u UserService) List(val *url.Values, claims *types.Claims) (types.Entity, *types.AppError) {
	if val == nil {
		return nil, newJSONError(errors.ErrMustProvideOptions, http.StatusInternalServerError)
	}
	var opt helpers.ListOptions
	if err := helpers.SchemaDecoder.Decode(&opt, *val); err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	// TODO: fix this
	users := make(models.Users, 0)
	sql := `SELECT id, email, 'password' AS password, name, role,
		created_at, updated_at, deleted_at
		FROM users
		WHERE verified IS TRUE
		AND deleted_at IS NULL;`
	if err := models.DBH.Select(&users, sql); err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}
	payload := payloads.Users{
		Users: &users,
		Meta: &models.UserMeta{
			CanAdd: claims.Role == "A",
		},
	}
	return &payload, nil
}

// Get retrieves a single user.
func (u UserService) Get(id int64, dummy string, claims *types.Claims) (types.Entity, *types.AppError) {
	user, err := models.DbGetUserByID(id)
	user.Password = ""
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	user.CanEdit = claims.Role == "A" || id == claims.Sub

	payload := payloads.User{
		User: user,
		Meta: &models.UserMeta{
			CanAdd: claims.Role == "A",
		},
	}
	return &payload, nil
}

// Update modifies an existing user.
func (u UserService) Update(id int64, e *types.Entity, dummy string, claims *types.Claims) *types.AppError {
	user := (*e).(*payloads.User).User

	originalUser, err := models.DbGetUserByID(id)
	if err != nil {
		return newJSONError(err, http.StatusInternalServerError)
	}

	user.ID = id
	user.Password = originalUser.Password
	user.Verified = originalUser.Verified
	user.UpdatedAt = helpers.CurrentTime()

	if err := user.Validate(); err != nil {
		return &types.AppError{Error: err, Status: helpers.StatusUnprocessableEntity}
	}

	// TODO: fix this
	count, err := models.DBH.Update(user)
	user.Password = ""
	if err != nil {
		return newJSONError(err, http.StatusInternalServerError)
	}
	if count != 1 {
		return newJSONError(errors.ErrUserNotUpdated, http.StatusInternalServerError)
	}

	return nil
}

// Create initializes a new user.
func (u UserService) Create(e *types.Entity, dummy string, claims *types.Claims) *types.AppError {
	user := (*e).(*payloads.User).User
	if err := user.Validate(); err != nil {
		return &types.AppError{Error: err, Status: helpers.StatusUnprocessableEntity}
	}
	ct := helpers.CurrentTime()
	user.CreatedAt = ct
	user.UpdatedAt = ct
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return newJSONError(err, http.StatusInternalServerError)
	}
	user.Password = string(hash)
	user.Role = "R"
	user.Verified = false

	// TODO: fix this
	if err := models.DBH.Insert(user.UserBase); err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code == "23505" {
				return newJSONError(errors.ErrEmailAddressTaken, http.StatusInternalServerError)
			}
		}
		return newJSONError(err, http.StatusInternalServerError)
	}

	user.Password = "password" // don't want to send the hashed PW back to the client

	q := `INSERT INTO verification (user_id, nonce, referer, created_at) VALUES ($1, $2, $3, $4);`
	// TODO: move helpers.GenerateNonce
	nonce, err := helpers.GenerateNonce()
	if err != nil {
		return newJSONError(err, http.StatusInternalServerError)
	}
	// TODO: fix this
	_, err = models.DBH.Exec(q, user.ID, nonce, claims.Ref, ct)
	if err != nil {
		return newJSONError(err, http.StatusInternalServerError)
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
			return newJSONError(err, http.StatusInternalServerError)
		}
	}

	return nil
}

// HandleUserVerify is a HTTP handler for verifiying a user.
func HandleUserVerify(w http.ResponseWriter, r *http.Request) *types.AppError {
	// TODO: clean this up
	nonce := mux.Vars(r)["Nonce"]
	q := `SELECT user_id, referer FROM verification WHERE nonce=$1;`

	var ver struct {
		UserID  int64
		Referer string
	}
	if err := models.DBH.SelectOne(&ver, q, nonce); err != nil {
		log.Print(err)
		return newJSONError(err, http.StatusInternalServerError)
	}

	if ver.UserID == 0 {
		return newJSONError(errors.ErrUserNotFound, http.StatusInternalServerError)
	}

	var user models.User
	if err := models.DBH.Get(&user, ver.UserID); err != nil {
		return newJSONError(err, http.StatusInternalServerError)
	}

	user.UpdatedAt = helpers.CurrentTime()
	user.Verified = true

	count, err := models.DBH.Update(&user)
	if err != nil {
		return newJSONError(err, http.StatusInternalServerError)
	}
	if count != 1 {
		return newJSONError(errors.ErrUserNotUpdated, http.StatusInternalServerError)
	}

	q = `DELETE FROM verification WHERE user_id=$1;`
	_, err = models.DBH.Exec(q, user.ID)
	if err != nil {
		return newJSONError(err, http.StatusInternalServerError)
	}
	fmt.Fprintln(w, `{"msg":"All set! Please log in."}`)
	return nil
}

// HandleUserLockout is a HTTP handler for unlocking a user's account.
func HandleUserLockout(w http.ResponseWriter, r *http.Request) *types.AppError {
	email := r.FormValue("email")
	if email == "" {
		return newJSONError(errors.ErrUserMissingEmail, http.StatusInternalServerError)
	}
	token, err := auth.Middleware.CreateToken(email)
	if err != nil {
		return newJSONError(err, http.StatusInternalServerError)
	}
	origin := r.Header.Get("Origin")
	hostURL, err := url.Parse(origin)
	if err != nil {
		return newJSONError(err, http.StatusInternalServerError)
	}
	hostURL.Path += "/users/lockoutauthenticate"
	params := url.Values{}
	params.Add("token", token)
	hostURL.RawQuery = params.Encode()

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
			mg.Domain(), hostURL.String())
		m := mailgun.NewMessage(sender, subject, message, recipient)
		_, _, err := mg.Send(m)
		if err != nil {
			log.Printf("%+v\n", err)
			return newJSONError(err, http.StatusInternalServerError)
		}
	}

	fmt.Fprintln(w, `{}`)
	return nil
}

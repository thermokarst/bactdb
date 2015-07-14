package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"github.com/mailgun/mailgun-go"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound           = errors.New("User not found")
	ErrUserNotFoundJSON       = newJSONError(ErrUserNotFound, http.StatusNotFound)
	ErrUserNotUpdated         = errors.New("User not updated")
	ErrUserNotUpdatedJSON     = newJSONError(ErrUserNotUpdated, http.StatusBadRequest)
	ErrInvalidEmailOrPassword = errors.New("Invalid email or password")
	ErrEmailAddressTaken      = errors.New("Email address already registered")
	ErrEmailAddressTakenJSON  = newJSONError(ErrEmailAddressTaken, http.StatusBadRequest)
)

func init() {
	DB.AddTableWithName(User{}, "users").SetKeys(true, "Id")
}

type UserService struct{}

type User struct {
	Id        int64    `json:"id,omitempty"`
	Email     string   `db:"email" json:"email"`
	Password  string   `db:"password" json:"password,omitempty"`
	Name      string   `db:"name" json:"name"`
	Role      string   `db:"role" json:"role"`
	Verified  bool     `db:"verified" json:"-"`
	CreatedAt NullTime `db:"created_at" json:"createdAt"`
	UpdatedAt NullTime `db:"updated_at" json:"updatedAt"`
	DeletedAt NullTime `db:"deleted_at" json:"deletedAt"`
}

type UserValidation struct {
	Email    []string `json:"email,omitempty"`
	Password []string `json:"password,omitempty"`
	Name     []string `json:"name,omitempty"`
	Role     []string `json:"role,omitempty"`
}

func (uv UserValidation) Error() string {
	errs, err := json.Marshal(struct {
		UserValidation `json:"errors"`
	}{uv})
	if err != nil {
		return err.Error()
	}
	return string(errs)
}

type Users []*User

type UserJSON struct {
	User *User `json:"user"`
}

type UsersJSON struct {
	Users *Users `json:"users"`
}

func (u *User) marshal() ([]byte, error) {
	return json.Marshal(&UserJSON{User: u})
}

func (u *Users) marshal() ([]byte, error) {
	return json.Marshal(&UsersJSON{Users: u})
}

func (u UserService) unmarshal(b []byte) (entity, error) {
	var uj UserJSON
	err := json.Unmarshal(b, &uj)
	return uj.User, err
}

func (u *User) validate() error {
	var uv UserValidation
	validationError := false

	if u.Name == "" {
		uv.Name = append(uv.Name, MustProvideAValue)
		validationError = true
	}

	if u.Email == "" {
		uv.Email = append(uv.Email, MustProvideAValue)
		validationError = true
	}

	regex, _ := regexp.Compile(`(\w[-._\w]*\w@\w[-._\w]*\w\.\w{2,3})`)
	if u.Email != "" && !regex.MatchString(u.Email) {
		uv.Email = append(uv.Email, "Must provide a valid email address")
		validationError = true
	}

	if u.Password == "" {
		uv.Password = append(uv.Password, MustProvideAValue)
		validationError = true
	} else {
		sevenOrMore, number, upper := verifyPassword(u.Password)
		if !sevenOrMore || !number || !upper {
			uv.Password = append(uv.Password, "Password must be at least 8 characters"+
				" long, and have at least one number and one uppercase letter")
			validationError = true
		}
	}

	if validationError {
		return uv
	}
	return nil
}

func (u UserService) list(val *url.Values, claims *Claims) (entity, *appError) {
	if val == nil {
		return nil, ErrMustProvideOptionsJSON
	}
	var opt ListOptions
	if err := schemaDecoder.Decode(&opt, *val); err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	users := make(Users, 0)
	sql := `SELECT id, email, 'password' AS password, name, role,
		created_at, updated_at, deleted_at
		FROM users
		WHERE verified IS TRUE
		AND deleted_at IS NULL;`
	if err := DBH.Select(&users, sql); err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}
	return &users, nil
}

func (u UserService) get(id int64, dummy string, claims *Claims) (entity, *appError) {
	var user User
	q := `SELECT id, email, 'password' AS password, name, role,
		created_at, updated_at, deleted_at
		FROM users
		WHERE id=$1
		AND verified IS TRUE
		AND deleted_at IS NULL;`
	if err := DBH.SelectOne(&user, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFoundJSON
		}
		return nil, newJSONError(err, http.StatusInternalServerError)
	}
	return &user, nil
}

func (u UserService) update(id int64, e *entity, dummy string, claims *Claims) *appError {
	user := (*e).(*User)
	user.UpdatedAt = currentTime()
	user.Id = id

	count, err := DBH.Update(user)
	if err != nil {
		return newJSONError(err, http.StatusInternalServerError)
	}
	if count != 1 {
		return ErrUserNotUpdatedJSON
	}
	return nil
}

func (u UserService) create(e *entity, dummy string, claims *Claims) *appError {
	user := (*e).(*User)
	if err := user.validate(); err != nil {
		return &appError{Error: err, Status: StatusUnprocessableEntity}
	}
	ct := currentTime()
	user.CreatedAt = ct
	user.UpdatedAt = ct
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return newJSONError(err, http.StatusInternalServerError)
	}
	user.Password = string(hash)
	user.Role = "R"
	user.Verified = false

	if err := DBH.Insert(user); err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code == "23505" {
				return ErrEmailAddressTakenJSON
			}
		}
		return newJSONError(err, http.StatusInternalServerError)
	}

	user.Password = "password" // don't want to send the hashed PW back to the client

	q := `INSERT INTO verification (user_id, nonce, referer, created_at) VALUES ($1, $2, $3, $4);`
	nonce, err := generateNonce()
	if err != nil {
		return newJSONError(err, http.StatusInternalServerError)
	}
	_, err = DBH.Exec(q, user.Id, nonce, claims.Ref, ct)
	if err != nil {
		return newJSONError(err, http.StatusInternalServerError)
	}

	// Send out confirmation email
	mg, ok := mgAccts[claims.Ref]
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

// for thermokarst/jwt: authentication callback
func dbAuthenticate(email string, password string) error {
	var user User
	q := `SELECT *
		FROM users
		WHERE lower(email)=lower($1)
		AND verified IS TRUE
		AND deleted_at IS NULL;`
	if err := DBH.SelectOne(&user, q, email); err != nil {
		return ErrInvalidEmailOrPassword
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return ErrInvalidEmailOrPassword
	}
	return nil
}

// for thermokarst/jwt: setting user in claims bundle
func dbGetUserByEmail(email string) (*User, error) {
	var user User
	q := `SELECT *
		FROM users
		WHERE lower(email)=lower($1)
		AND verified IS TRUE
		AND deleted_at IS NULL;`
	if err := DBH.SelectOne(&user, q, email); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func handleUserVerify(w http.ResponseWriter, r *http.Request) {
	nonce := mux.Vars(r)["Nonce"]
	q := `SELECT user_id, referer FROM verification WHERE nonce=$1;`

	var ver struct {
		User_id int64
		Referer string
	}
	if err := DBH.SelectOne(&ver, q, nonce); err != nil {
		log.Print(err)
		fmt.Fprintln(w, "Invalid URL")
		return
	}

	failURL := fmt.Sprintf("%s/users/new/fail", ver.Referer)
	successURL := fmt.Sprintf("%s/users/new/success", ver.Referer)

	if ver.User_id == 0 {
		http.Redirect(w, r, failURL, http.StatusMovedPermanently)
		return
	}

	var user User
	if err := DBH.Get(&user, ver.User_id); err != nil {
		http.Redirect(w, r, failURL, http.StatusMovedPermanently)
		return
	}

	user.UpdatedAt = currentTime()
	user.Verified = true

	count, err := DBH.Update(&user)
	if err != nil {
		http.Redirect(w, r, failURL, http.StatusMovedPermanently)
		return
	}
	if count != 1 {
		http.Redirect(w, r, failURL, http.StatusMovedPermanently)
		return
	}

	q = `DELETE FROM verification WHERE user_id=$1;`
	_, err = DBH.Exec(q, user.Id)
	if err != nil {
		http.Redirect(w, r, failURL, http.StatusMovedPermanently)
		return
	}
	http.Redirect(w, r, successURL, http.StatusMovedPermanently)
}

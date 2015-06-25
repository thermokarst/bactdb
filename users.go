package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound           = errors.New("User not found")
	ErrUserNotFoundJSON       = newJSONError(ErrUserNotFound, http.StatusNotFound)
	ErrInvalidEmailOrPassword = errors.New("Invalid email or password")
)

func init() {
	DB.AddTableWithName(User{}, "users").SetKeys(true, "Id")
}

type UserService struct{}

type User struct {
	Id        int64     `json:"id,omitempty"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password" json:"-"`
	Name      string    `db:"name" json:"name"`
	Role      string    `db:"role" json:"role"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
	DeletedAt NullTime  `db:"deleted_at" json:"deletedAt"`
}

type UserValidation struct {
	Email    []string `json:"email,omitempty"`
	Password []string `json:"password,omitempty"`
	Name     []string `json:"name,omitempty"`
	Role     []string `json:"role,omitempty"`
}

func (uv UserValidation) Error() string {
	errs, err := json.Marshal(uv)
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
		uv.Name = append(uv.Name, "Must provide a value")
		validationError = true
	}

	if validationError {
		errs, _ := json.Marshal(uv)
		return errors.New(string(errs))
	}
	return nil
}

func (u UserService) list(val *url.Values) (entity, error) {
	if val == nil {
		return nil, errors.New("must provide options")
	}
	var opt ListOptions
	if err := schemaDecoder.Decode(&opt, *val); err != nil {
		return nil, err
	}

	users := make(Users, 0)
	sql := `SELECT * FROM users;`
	if err := DBH.Select(&users, sql); err != nil {
		return nil, err
	}
	return &users, nil
}

func (u UserService) get(id int64, genus string) (entity, *appError) {
	var user User
	q := `SELECT * FROM users WHERE id=$1;`
	if err := DBH.SelectOne(&user, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFoundJSON
		}
		return nil, newJSONError(err, http.StatusInternalServerError)
	}
	return &user, nil
}

func (u UserService) update(id int64, e *entity, claims Claims) error {
	user := (*e).(*User)
	user.UpdatedAt = time.Now()
	user.Id = id

	count, err := DBH.Update(user)
	if err != nil {
		return err
	}
	if count != 1 {
		return ErrStrainNotUpdated
	}
	return nil
}

func (u UserService) create(e *entity, claims Claims) error {
	user := (*e).(*User)
	if err := user.validate(); err != nil {
		return err
	}
	ct := time.Now()
	user.CreatedAt = ct
	user.UpdatedAt = ct
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return err
	}
	user.Password = string(hash)

	if err := DBH.Insert(user); err != nil {
		return err
	}
	return nil
}

// for thermokarst/jwt: authentication callback
func dbAuthenticate(email string, password string) error {
	var user User
	q := `SELECT * FROM users WHERE lower(email)=lower($1);`
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
	q := `SELECT * FROM users WHERE lower(email)=lower($1);`
	if err := DBH.SelectOne(&user, q, email); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

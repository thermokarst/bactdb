package models

import (
	"database/sql"
	"encoding/json"
	"regexp"

	"github.com/thermokarst/bactdb/Godeps/_workspace/src/golang.org/x/crypto/bcrypt"
	"github.com/thermokarst/bactdb/errors"
	"github.com/thermokarst/bactdb/helpers"
	"github.com/thermokarst/bactdb/types"
)

func init() {
	DB.AddTableWithName(UserBase{}, "users").SetKeys(true, "Id")
}

type UserBase struct {
	Id        int64          `json:"id,omitempty"`
	Email     string         `db:"email" json:"email"`
	Password  string         `db:"password" json:"password,omitempty"`
	Name      string         `db:"name" json:"name"`
	Role      string         `db:"role" json:"role"`
	Verified  bool           `db:"verified" json:"-"`
	CreatedAt types.NullTime `db:"created_at" json:"createdAt"`
	UpdatedAt types.NullTime `db:"updated_at" json:"updatedAt"`
	DeletedAt types.NullTime `db:"deleted_at" json:"deletedAt"`
}

type User struct {
	*UserBase
	CanEdit bool `db:"-" json:"canEdit"`
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

type UserMeta struct {
	CanAdd bool `json:"canAdd"`
}

func (u *Users) Marshal() ([]byte, error) {
	return json.Marshal(&UsersJSON{Users: u})
}

func (u *User) Validate() error {
	var uv UserValidation
	validationError := false

	if u.Name == "" {
		uv.Name = append(uv.Name, helpers.MustProvideAValue)
		validationError = true
	}

	if u.Email == "" {
		uv.Email = append(uv.Email, helpers.MustProvideAValue)
		validationError = true
	}

	regex, _ := regexp.Compile(`(\w[-._\w]*\w@\w[-._\w]*\w\.\w{2,3})`)
	if u.Email != "" && !regex.MatchString(u.Email) {
		uv.Email = append(uv.Email, "Must provide a valid email address")
		validationError = true
	}

	if len(u.Password) < 8 {
		uv.Password = append(uv.Password, "Password must be at least 8 characters")
		validationError = true
	}

	if validationError {
		return uv
	}
	return nil
}

// for thermokarst/jwt: authentication callback
func DbAuthenticate(email string, password string) error {
	var user User
	q := `SELECT *
		FROM users
		WHERE lower(email)=lower($1)
		AND verified IS TRUE
		AND deleted_at IS NULL;`
	if err := DBH.SelectOne(&user, q, email); err != nil {
		return errors.InvalidEmailOrPassword
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return errors.InvalidEmailOrPassword
	}
	return nil
}

func DbGetUserById(id int64) (*User, error) {
	var user User
	q := `SELECT *
		FROM users
		WHERE id=$1
		AND verified IS TRUE
		AND deleted_at IS NULL;`
	if err := DBH.SelectOne(&user, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.UserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// for thermokarst/jwt: setting user in claims bundle
func DbGetUserByEmail(email string) (*User, error) {
	var user User
	q := `SELECT *
		FROM users
		WHERE lower(email)=lower($1)
		AND verified IS TRUE
		AND deleted_at IS NULL;`
	if err := DBH.SelectOne(&user, q, email); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.UserNotFound
		}
		return nil, err
	}
	return &user, nil
}

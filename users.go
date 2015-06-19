package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/url"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound           = errors.New("user not found")
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

func (u UserService) get(id int64, genus string) (entity, error) {
	var user User
	q := `SELECT * FROM users WHERE id=$1;`
	if err := DBH.SelectOne(&user, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
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

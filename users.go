package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound           = errors.New("user not found")
	ErrInvalidEmailOrPassword = errors.New("Invalid email or password")
)

func init() {
	DB.AddTableWithName(User{}, "users").SetKeys(true, "Id")
}

// A User is a person that has administrative access to bactdb.
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

type UserJSON struct {
	User *User `json:"user"`
}

type UsersJSON struct {
	Users []*User `json:"users"`
}

func serveUsersList(w http.ResponseWriter, r *http.Request) {
	var opt ListOptions
	if err := schemaDecoder.Decode(&opt, r.URL.Query()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	users, err := dbGetUsers(&opt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if users == nil {
		users = []*User{}
	}
	data, err := json.Marshal(UsersJSON{Users: users})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(data)
}

func serveUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["Id"], 10, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := dbGetUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(UserJSON{User: user})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(data)
}

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

func dbGetUsers(opt *ListOptions) ([]*User, error) {
	var users []*User
	sql := `SELECT * FROM users;`
	if err := DBH.Select(&users, sql); err != nil {
		return nil, err
	}
	return users, nil
}

func dbGetUser(id int64) (*User, error) {
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

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

func init() {
	DB.AddTableWithName(User{}, "users").SetKeys(true, "Id")
}

// A User is a person that has administrative access to bactdb.
// Todo: add password
type User struct {
	Id        int64     `json:"id,omitempty"`
	Username  string    `db:"username" json:"username"`
	Password  string    `db:"password" json:"-"`
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

func (m *User) String() string {
	return fmt.Sprintf("%v", *m)
}

type UserSession struct {
	Token       string `json:"token"`
	AccessLevel string `json:"access_level"`
	Genus       string `json:"genus"`
}

func serveAuthenticateUser(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	user_session, err := dbAuthenticate(username, password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	t := jwt.New(jwt.GetSigningMethod("RS256"))
	t.Claims["auth_level"] = user_session.AccessLevel
	t.Claims["genus"] = user_session.Genus
	t.Claims["exp"] = time.Now().Add(time.Minute * 60 * 24).Unix()
	tokenString, err := t.SignedString(signKey)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user_session.Token = tokenString
	data, err := json.Marshal(user_session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(data)
}

func dbAuthenticate(username string, password string) (*UserSession, error) {
	var users []User
	var user_session UserSession

	if err := DBH.Select(&users, `SELECT * FROM users WHERE username=$1;`, username); err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, ErrUserNotFound
	}
	if err := bcrypt.CompareHashAndPassword([]byte(users[0].Password), []byte(password)); err != nil {
		return nil, err
	}
	user_session.AccessLevel = "read"
	user_session.Genus = "hymenobacter"
	return &user_session, nil
}

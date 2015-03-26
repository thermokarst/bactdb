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
	*User
	Role  string `json:"access_level"`
	Genus string `json:"genus"`
}

func serveAuthenticateUser(w http.ResponseWriter, r *http.Request) {
	var a struct {
		Username string
		Password string
	}
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	user_session, err := dbAuthenticate(a.Username, a.Password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"Invalid username or password"}`))
		return
	}

	currentTime := time.Now()

	t := jwt.New(jwt.GetSigningMethod("HS256"))
	t.Claims["name"] = user_session.Username
	t.Claims["iss"] = "bactdb"
	t.Claims["sub"] = "user@example.com" // TODO: fix this
	t.Claims["role"] = user_session.Role
	t.Claims["iat"] = currentTime.Unix()
	t.Claims["exp"] = currentTime.Add(time.Minute * 60 * 24).Unix()
	tokenString, err := t.SignedString(signKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var token struct {
		Token string `json:"token"`
	}
	token.Token = tokenString
	data, err := json.Marshal(token)
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
	user_session.User = &users[0]
	user_session.Role = "admin"
	user_session.Genus = "hymenobacter"
	return &user_session, nil
}

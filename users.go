package main

import (
	"encoding/json"
	"errors"
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

func serveAuthenticateUser(w http.ResponseWriter, r *http.Request) {
	var a struct {
		Email    string
		Password string
	}
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	user_session, err := dbAuthenticate(a.Email, a.Password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"Invalid email or password"}`))
		return
	}

	currentTime := time.Now()

	t := jwt.New(jwt.SigningMethodRS256)
	t.Claims["name"] = user_session.Name
	t.Claims["iss"] = "bactdb"
	t.Claims["sub"] = user_session.Email
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

func dbAuthenticate(email string, password string) (*User, error) {
	var users []User
	if err := DBH.Select(&users, `SELECT * FROM users WHERE lower(email)=lower($1);`, email); err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, ErrUserNotFound
	}
	if err := bcrypt.CompareHashAndPassword([]byte(users[0].Password), []byte(password)); err != nil {
		return nil, err
	}
	return &users[0], nil
}

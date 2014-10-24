package models

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/thermokarst/bactdb/router"
)

// A User is a person that has administrative access to bactdb.
type User struct {
	Id        int64     `json:"id,omitempty"`
	UserName  string    `json:"user_name"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	DeletedAt time.Time `db:"deleted_at" json:"deleted_at"`
}

func NewUser() *User {
	return &User{UserName: "Test User"}
}

// UsersService interacts with the user-related endpoints in bactdb's API.
type UsersService interface {
	// Get a user.
	Get(id int64) (*User, error)

	// List all users.
	List(opt *UserListOptions) ([]*User, error)

	// Create a new user. The newly created user's ID is written to user.Id
	Create(user *User) (created bool, err error)
}

var (
	ErrUserNotFound = errors.New("user not found")
)

type usersService struct {
	client *Client
}

func (s *usersService) Get(id int64) (*User, error) {
	// Pass in key value pairs as strings, so that the gorilla mux URL
	// generation is happy.
	strId := strconv.FormatInt(id, 10)

	url, err := s.client.url(router.User, map[string]string{"Id": strId}, nil)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	var user *User
	_, err = s.client.Do(req, &user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *usersService) Create(user *User) (bool, error) {
	url, err := s.client.url(router.CreateUser, nil, nil)
	if err != nil {
		return false, err
	}

	req, err := s.client.NewRequest("POST", url.String(), user)
	if err != nil {
		return false, err
	}

	resp, err := s.client.Do(req, &user)
	if err != nil {
		return false, err
	}

	return resp.StatusCode == http.StatusCreated, nil
}

type UserListOptions struct {
	ListOptions
}

func (s *usersService) List(opt *UserListOptions) ([]*User, error) {
	url, err := s.client.url(router.Users, nil, opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	var users []*User
	_, err = s.client.Do(req, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

type MockUsersService struct {
	Get_    func(id int64) (*User, error)
	List_   func(opt *UserListOptions) ([]*User, error)
	Create_ func(user *User) (bool, error)
}

var _ UsersService = &MockUsersService{}

func (s *MockUsersService) Get(id int64) (*User, error) {
	if s.Get_ == nil {
		return nil, nil
	}
	return s.Get_(id)
}

func (s *MockUsersService) Create(user *User) (bool, error) {
	if s.Create_ == nil {
		return false, nil
	}
	return s.Create_(user)
}

func (s *MockUsersService) List(opt *UserListOptions) ([]*User, error) {
	if s.List_ == nil {
		return nil, nil
	}
	return s.List_(opt)
}

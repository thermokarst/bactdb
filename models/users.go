package models

import (
	"errors"
	"strconv"
	"time"

	"github.com/thermokarst/bactdb/router"
)

// A User is a person that has administrative access to bactdb.
type User struct {
	Id        int64     `json:"id"`
	UserName  string    `sql:"size:100" json:"user_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

// UsersService interacts with the user-related endpoints in bactdb's API.
type UsersService interface {
	// Get a user.
	Get(id int64) (*User, error)

	// List all users.
	List(opt *UserListOptions) ([]*User, error)
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
	Get_  func(id int64) (*User, error)
	List_ func(opt *UserListOptions) ([]*User, error)
}

func (s *MockUsersService) Get(id int64) (*User, error) {
	if s.Get_ == nil {
		return nil, nil
	}
	return s.Get_(id)
}

func (s *MockUsersService) List(opt *UserListOptions) ([]*User, error) {
	if s.List_ == nil {
		return nil, nil
	}
	return s.List_(opt)
}

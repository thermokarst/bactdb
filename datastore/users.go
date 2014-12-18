package datastore

import (
	"fmt"
	"strings"
	"time"

	"github.com/thermokarst/bactdb/models"
)

func init() {
	DB.AddTableWithName(models.User{}, "users").SetKeys(true, "Id")
}

type usersStore struct {
	*Datastore
}

func (s *usersStore) Get(id int64) (*models.User, error) {
	var users []*models.User
	if err := s.dbh.Select(&users, `SELECT * FROM users WHERE id=$1;`, id); err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, models.ErrUserNotFound
	}
	return users[0], nil
}

func (s *usersStore) Create(user *models.User) (bool, error) {
	currentTime := time.Now()
	user.CreatedAt = currentTime
	user.UpdatedAt = currentTime
	fmt.Println(user)
	if err := s.dbh.Insert(user); err != nil {
		if strings.Contains(err.Error(), `violates unique constraint "username_idx"`) {
			return false, err
		}
	}
	return true, nil
}

func (s *usersStore) List(opt *models.UserListOptions) ([]*models.User, error) {
	if opt == nil {
		opt = &models.UserListOptions{}
	}
	var users []*models.User
	err := s.dbh.Select(&users, `SELECT * FROM users LIMIT $1 OFFSET $2;`, opt.PerPageOrDefault(), opt.Offset())
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *usersStore) Authenticate(username string, password string) (*string, error) {
	var users []*models.User
	if err := s.dbh.Select(&users, `SELECT * FROM users WHERE username=$1;`, username); err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, models.ErrUserNotFound
	}
	auth_level := "read"
	return &auth_level, nil
}

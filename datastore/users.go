package datastore

import (
	"strings"
	"time"

	"github.com/thermokarst/bactdb/models"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	DB.AddTableWithName(models.User{}, "users").SetKeys(true, "Id")
}

type usersStore struct {
	*Datastore
}

func (s *usersStore) Get(id int64) (*models.User, error) {
	var user models.User
	if err := s.dbh.SelectOne(&user, `SELECT * FROM users WHERE id=$1;`, id); err != nil {
		return nil, err
	}
	if &user == nil {
		return nil, models.ErrUserNotFound
	}
	return &user, nil
}

func (s *usersStore) Create(user *models.User) (bool, error) {
	currentTime := time.Now()
	user.CreatedAt = currentTime
	user.UpdatedAt = currentTime
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		panic(err)
	}
	user.Password = string(hashedPassword)
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

func (s *usersStore) Authenticate(username string, password string) (*models.UserSession, error) {
	var users []*models.User
	var user_session models.UserSession

	if err := s.dbh.Select(&users, `SELECT * FROM users WHERE username=$1;`, username); err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, models.ErrUserNotFound
	}
	if err := bcrypt.CompareHashAndPassword([]byte(users[0].Password), []byte(password)); err != nil {
		return nil, err
	}
	user_session.AccessLevel = "read"
	user_session.Genus = "hymenobacter"
	return &user_session, nil
}

package datastore

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/jmoiron/modl"
	"github.com/thermokarst/bactdb/models"
)

func init() {
	DB.AddTableWithName(models.User{}, "users").SetKeys(true, "Id")
	createSQL = append(createSQL,
		`CREATE UNIQUE INDEX username_idx ON users (username);`,
	)
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
	retries := 3
	var wantRetry bool

retry:
	retries--
	wantRetry = false
	if retries == 0 {
		return false, fmt.Errorf("failed to create user with username %q after retrying", user.UserName)
	}

	var created bool
	err := transact(s.dbh, func(tx modl.SqlExecutor) error {
		var existing []*models.User
		if err := tx.Select(&existing, `SELECT * FROM users WHERE username=$1 LIMIT 1;`, user.UserName); err != nil {
			return err
		}
		if len(existing) > 0 {
			*user = *existing[0]
			return nil
		}

		if err := tx.Insert(user); err != nil {
			if strings.Contains(err.Error(), `violates unique constraint "username_idx"`) {
				time.Sleep(time.Duration(rand.Intn(75)) * time.Millisecond)
				wantRetry = true
				return err
			}
			return err
		}

		created = true
		return nil
	})
	if wantRetry {
		goto retry
	}
	return created, err
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

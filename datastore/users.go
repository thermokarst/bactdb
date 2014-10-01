package datastore

import "github.com/thermokarst/bactdb/models"

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

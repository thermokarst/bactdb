package models

import (
	"database/sql"
	"regexp"

	"github.com/thermokarst/bactdb/Godeps/_workspace/src/github.com/jmoiron/modl"
	"github.com/thermokarst/bactdb/Godeps/_workspace/src/golang.org/x/crypto/bcrypt"
	"github.com/thermokarst/bactdb/errors"
	"github.com/thermokarst/bactdb/helpers"
	"github.com/thermokarst/bactdb/types"
)

func init() {
	DB.AddTableWithName(UserBase{}, "users").SetKeys(true, "ID")
}

// PreInsert is a modl hook.
func (u *UserBase) PreInsert(e modl.SqlExecutor) error {
	ct := helpers.CurrentTime()
	u.CreatedAt = ct
	u.UpdatedAt = ct
	return nil
}

// PreUpdate is a modl hook.
func (u *UserBase) PreUpdate(e modl.SqlExecutor) error {
	u.UpdatedAt = helpers.CurrentTime()
	return nil
}

// UpdateError satisfies base interface.
func (u *UserBase) UpdateError() error {
	return errors.ErrUserNotUpdated
}

// DeleteError satisfies base interface.
func (u *UserBase) DeleteError() error {
	return errors.ErrUserNotDeleted
}

func (u *UserBase) validate() types.ValidationError {
	uv := make(types.ValidationError, 0)

	if u.Name == "" {
		uv = append(uv, types.NewValidationError(
			"name",
			helpers.MustProvideAValue))
	}

	if u.Email == "" {
		uv = append(uv, types.NewValidationError(
			"email",
			helpers.MustProvideAValue))
	}

	regex, _ := regexp.Compile(`(\w[-._\w]*\w@\w[-._\w]*\w\.\w{2,3})`)
	if u.Email != "" && !regex.MatchString(u.Email) {
		uv = append(uv, types.NewValidationError(
			"email",
			"Must provide a valid email address"))
	}

	if len(u.Password) < 8 {
		uv = append(uv, types.NewValidationError(
			"password",
			"Password must be at least 8 characters"))
	}

	if len(uv) > 0 {
		return uv
	}

	return nil
}

// UserBase is what the DB expects to see for write operations.
type UserBase struct {
	ID        int64          `json:"id,omitempty"`
	Email     string         `db:"email" json:"email"`
	Password  string         `db:"password" json:"password,omitempty"`
	Name      string         `db:"name" json:"name"`
	Role      string         `db:"role" json:"role"`
	Verified  bool           `db:"verified" json:"-"`
	CreatedAt types.NullTime `db:"created_at" json:"createdAt"`
	UpdatedAt types.NullTime `db:"updated_at" json:"updatedAt"`
}

// User is what the DB expects to see for read operations, and is what the API
// expects to return to the requester.
type User struct {
	*UserBase
	CanEdit bool `db:"-" json:"canEdit"`
}

// UserValidation handles validation of a user record.
type UserValidation struct {
	Email    []string `json:"email,omitempty"`
	Password []string `json:"password,omitempty"`
	Name     []string `json:"name,omitempty"`
	Role     []string `json:"role,omitempty"`
}

// Users are multiple user entities.
type Users []*User

// UserMeta stashes some metadata related to the entity.
type UserMeta struct {
	CanAdd bool `json:"canAdd"`
}

// DbAuthenticate authenticates a user.
// For thermokarst/jwt: authentication callback
func DbAuthenticate(email string, password string) error {
	var user User
	q := `SELECT *
		FROM users
		WHERE lower(email)=lower($1)
		AND verified IS TRUE;`
	if err := DBH.SelectOne(&user, q, email); err != nil {
		return errors.ErrInvalidEmailOrPassword
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return errors.ErrInvalidEmailOrPassword
	}
	return nil
}

// GetUser returns a specific user record by ID.
func GetUser(id int64, dummy string, claims *types.Claims) (*User, error) {
	var user User
	q := `SELECT *
		FROM users
		WHERE id=$1
		AND verified IS TRUE;`
	if err := DBH.SelectOne(&user, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrUserNotFound
		}
		return nil, err
	}

	user.CanEdit = claims.Role == "A" || id == claims.Sub

	return &user, nil
}

// DbGetUserByEmail returns a specific user record by email.
// For thermokarst/jwt: setting user in claims bundle
func DbGetUserByEmail(email string) (*User, error) {
	var user User
	q := `SELECT *
		FROM users
		WHERE lower(email)=lower($1)
		AND verified IS TRUE;`
	if err := DBH.SelectOne(&user, q, email); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// ListUsers returns all users.
func ListUsers(opt helpers.ListOptions, claims *types.Claims) (*Users, error) {
	q := `SELECT id, email, 'password' AS password, name, role, created_at, updated_at
		FROM users
		WHERE verified IS TRUE;`

	users := make(Users, 0)
	if err := DBH.Select(&users, q); err != nil {
		return nil, err
	}

	for _, u := range users {
		u.CanEdit = claims.Role == "A" || u.ID == claims.Sub
	}

	return &users, nil
}

func UpdateUserPassword(claims *types.Claims, password string) error {
	user, err := GetUser(claims.Sub, "", claims)
	if err != nil {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	user.Password = string(hash)

	count, err := DBH.Update(user.UserBase)
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.ErrUserNotUpdated
	}
	return nil
}

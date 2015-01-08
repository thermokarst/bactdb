package datastore

import (
	"reflect"
	"testing"

	"github.com/jmoiron/modl"
	"github.com/thermokarst/bactdb/models"
	"golang.org/x/crypto/bcrypt"
)

func insertUser(t *testing.T, tx *modl.Transaction) *models.User {
	// Test on a clean database
	tx.Exec(`DELETE FROM users;`)

	user := newUser()
	if err := tx.Insert(user); err != nil {
		t.Fatal(err)
	}
	return user
}

func newUser() *models.User {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), 10)
	return &models.User{
		Username: "Test User",
		Password: string(hashedPassword),
	}
}

func TestUsersStore_Get_db(t *testing.T) {
	tx, _ := DB.Begin()
	defer tx.Rollback()

	want := insertUser(t, tx)

	d := NewDatastore(tx)

	user, err := d.Users.Get(want.Id)
	if err != nil {
		t.Fatal(err)
	}

	normalizeTime(&want.CreatedAt, &want.UpdatedAt, &want.DeletedAt)
	normalizeTime(&user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
	if !reflect.DeepEqual(user, want) {
		t.Errorf("got user %+v, want %+v", user, want)
	}
}

func TestUsersStore_Create_db(t *testing.T) {
	tx, _ := DB.Begin()
	defer tx.Rollback()

	user := newUser()

	d := NewDatastore(tx)

	created, err := d.Users.Create(user)
	if err != nil {
		t.Fatal(err)
	}

	if !created {
		t.Error("!created")
	}
	if user.Id == 0 {
		t.Error("want nonzero user.Id after submitting")
	}
}

func TestUsersStore_List_db(t *testing.T) {
	tx, _ := DB.Begin()
	defer tx.Rollback()

	user := insertUser(t, tx)
	want := []*models.User{user}

	d := NewDatastore(tx)

	users, err := d.Users.List(&models.UserListOptions{ListOptions: models.ListOptions{Page: 1, PerPage: 10}})
	if err != nil {
		t.Fatal(err)
	}

	for i := range want {
		normalizeTime(&want[i].CreatedAt, &want[i].UpdatedAt, &want[i].DeletedAt)
		normalizeTime(&users[i].CreatedAt, &users[i].UpdatedAt, &users[i].DeletedAt)
	}
	if !reflect.DeepEqual(users, want) {
		t.Errorf("got users %+v, want %+v", users, want)
	}
}

func TestUsersStore_Authenticate_db(t *testing.T) {
	tx, _ := DB.Begin()
	defer tx.Rollback()

	user := insertUser(t, tx)

	want := &models.UserSession{
		AccessLevel: "read",
		Genus:       "hymenobacter",
	}

	d := NewDatastore(tx)

	user_session, err := d.Users.Authenticate(user.Username, "password")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(user_session, want) {
		t.Errorf("got session %+v, want %+v", user_session, want)
	}
}

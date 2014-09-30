package datastore

import (
	"reflect"
	"testing"

	"github.com/thermokarst/bactdb/models"
)

func TestUsersStore_Get_db(t *testing.T) {
	want := &models.User{Id: 1, UserName: "Test User"}

	tx, _ := DB.Begin()
	defer tx.Rollback()
	// Test on a clean database
	tx.Exec(`DELETE FROM users;`)
	if err := tx.Insert(want); err != nil {
		t.Fatal(err)
	}

	d := NewDatastore(tx)
	user, err := d.Users.Get(1)
	if err != nil {
		t.Fatal(err)
	}

	normalizeTime(&want.CreatedAt, &want.UpdatedAt, &want.DeletedAt)
	if !reflect.DeepEqual(user, want) {
		t.Errorf("got user %+v, want %+v", user, want)
	}
}

func TestUsersStore_List_db(t *testing.T) {
	want := []*models.User{{Id: 1, UserName: "Test User"}}

	// tx := DBH
	tx, _ := DB.Begin()
	defer tx.Rollback()

	// Test on a clean database
	tx.Exec(`DELETE FROM users;`)
	if err := tx.Insert(want[0]); err != nil {
		t.Fatal(err)
	}

	d := NewDatastore(tx)
	users, err := d.Users.List(&models.UserListOptions{ListOptions: models.ListOptions{Page: 1, PerPage: 10}})
	if err != nil {
		t.Fatal(err)
	}

	for _, u := range want {
		normalizeTime(&u.CreatedAt, &u.UpdatedAt, &u.DeletedAt)
	}
	if !reflect.DeepEqual(users, want) {
		t.Errorf("got users %+v, want %+v", users, want)
	}
}

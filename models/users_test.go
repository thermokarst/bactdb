package models

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/thermokarst/bactdb/router"
)

func newUser() *User {
	user := NewUser()
	user.Id = 1
	return user
}

func TestUsersService_Get(t *testing.T) {
	setup()
	defer teardown()

	want := newUser()

	var called bool
	mux.HandleFunc(urlPath(t, router.User, map[string]string{"Id": "1"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")

		writeJSON(w, want)
	})

	user, err := client.Users.Get(1)
	if err != nil {
		t.Errorf("Users.Get returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	normalizeTime(&want.CreatedAt, &want.UpdatedAt, &want.DeletedAt)

	if !reflect.DeepEqual(user, want) {
		t.Errorf("Users.Get returned %+v, want %+v", user, want)
	}
}

func TestUsersService_Create(t *testing.T) {
	setup()
	defer teardown()

	want := newUser()

	var called bool
	mux.HandleFunc(urlPath(t, router.CreateUser, nil), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "POST")
		testBody(t, r, `{"id":1,"user_name":"Test User","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","deleted_at":"0001-01-01T00:00:00Z"}`+"\n")

		w.WriteHeader(http.StatusCreated)
		writeJSON(w, want)
	})

	user := newUser()
	created, err := client.Users.Create(user)
	if err != nil {
		t.Errorf("Users.Create returned error: %v", err)
	}

	if !created {
		t.Error("!created")
	}

	if !called {
		t.Fatal("!called")
	}

	normalizeTime(&want.CreatedAt, &want.UpdatedAt, &want.DeletedAt)
	if !reflect.DeepEqual(user, want) {
		t.Errorf("Users.Create returned %+v, want %+v", user, want)
	}
}

func TestUsersService_List(t *testing.T) {
	setup()
	defer teardown()

	want := []*User{newUser()}

	var called bool
	mux.HandleFunc(urlPath(t, router.Users, nil), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")
		testFormValues(t, r, values{})

		writeJSON(w, want)
	})

	users, err := client.Users.List(nil)
	if err != nil {
		t.Errorf("Users.List returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	for _, u := range want {
		normalizeTime(&u.CreatedAt, &u.UpdatedAt, &u.DeletedAt)
	}
	if !reflect.DeepEqual(users, want) {
		t.Errorf("Users.List return %+v, want %+v", users, want)
	}
}

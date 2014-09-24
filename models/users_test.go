package models

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/thermokarst/bactdb/router"
)

func TestUsersService_Get(t *testing.T) {
	setup()
	defer teardown()

	want := &User{Id: 1, UserName: "Test User"}

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

package api

import (
	"testing"

	"github.com/thermokarst/bactdb/models"
)

func newUser() *models.User {
	user := models.NewUser()
	user.Id = 1
	return user
}

func TestUser_Get(t *testing.T) {
	setup()

	wantUser := newUser()

	calledGet := false
	store.Users.(*models.MockUsersService).Get_ = func(id int64) (*models.User, error) {
		if id != wantUser.Id {
			t.Errorf("wanted request for user %d but got %d", wantUser.Id, id)
		}
		calledGet = true
		return wantUser, nil
	}

	gotUser, err := apiClient.Users.Get(wantUser.Id)
	if err != nil {
		t.Fatal(err)
	}

	if !calledGet {
		t.Error("!calledGet")
	}
	if !normalizeDeepEqual(wantUser, gotUser) {
		t.Errorf("got user %+v but wanted user %+v", wantUser, gotUser)
	}
}

func TestUser_Create(t *testing.T) {
	setup()

	wantUser := newUser()

	calledPost := false
	store.Users.(*models.MockUsersService).Create_ = func(user *models.User) (bool, error) {
		if !normalizeDeepEqual(wantUser, user) {
			t.Errorf("wanted request for user %d but got %d", wantUser, user)
		}
		calledPost = true
		return true, nil
	}

	success, err := apiClient.Users.Create(wantUser)
	if err != nil {
		t.Fatal(err)
	}

	if !calledPost {
		t.Error("!calledPost")
	}
	if !success {
		t.Error("!success")
	}
}

func TestUser_List(t *testing.T) {
	setup()

	wantUsers := []*models.User{newUser()}
	wantOpt := &models.UserListOptions{ListOptions: models.ListOptions{Page: 1, PerPage: 10}}

	calledList := false
	store.Users.(*models.MockUsersService).List_ = func(opt *models.UserListOptions) ([]*models.User, error) {
		if !normalizeDeepEqual(wantOpt, opt) {
			t.Errorf("wanted options %d but got %d", wantOpt, opt)
		}
		calledList = true
		return wantUsers, nil
	}

	users, err := apiClient.Users.List(wantOpt)
	if err != nil {
		t.Fatal(err)
	}

	if !calledList {
		t.Error("!calledList")
	}

	if !normalizeDeepEqual(&wantUsers, &users) {
		t.Errorf("got users %+v but wanted users %+v", users, wantUsers)
	}
}

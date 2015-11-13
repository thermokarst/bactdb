package payloads

import (
	"encoding/json"

	"github.com/thermokarst/bactdb/models"
)

// User is a payload that sideloads all of the necessary entities for a
// particular user.
type User struct {
	User *models.User `json:"user"`
}

// Users is a payload that sideloads all of the necessary entities for
// multiple users.
type Users struct {
	Users *models.Users `json:"users"`
}

// Marshal satisfies the CRUD interfaces.
func (u *User) Marshal() ([]byte, error) {
	return json.Marshal(u)
}

// Marshal satisfies the CRUD interfaces.
func (u *Users) Marshal() ([]byte, error) {
	return json.Marshal(u)
}

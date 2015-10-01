package payloads

import (
	"encoding/json"

	"github.com/thermokarst/bactdb/models"
)

type User struct {
	User *models.User     `json:"user"`
	Meta *models.UserMeta `json:"meta"`
}

func (u *User) Marshal() ([]byte, error) {
	return json.Marshal(u)
}

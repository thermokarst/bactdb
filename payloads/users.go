package payloads

import (
	"encoding/json"

	"github.com/thermokarst/bactdb/models"
)

type UserPayload struct {
	User *models.User     `json:"user"`
	Meta *models.UserMeta `json:"meta"`
}

func (u *UserPayload) Marshal() ([]byte, error) {
	return json.Marshal(u)
}

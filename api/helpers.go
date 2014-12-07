package api

import (
	"encoding/json"
	"log"
	"net/http"
)

// writeJSON writes a JSON Content-Type header and a JSON-encoded object to
// the http.ResponseWriter.
func writeJSON(w http.ResponseWriter, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	w.Header().Set("content-type", "application/json; charset=utf-8")
	_, err = w.Write(data)
	return err
}

// Message is for returning simple message payloads to the user
type Message struct {
	Message string `json:"message"`
}

// Error is for returning simple error payloads to the user, as well as logging
type Error struct {
	Error error
}

func (e Error) MarshalJSON() ([]byte, error) {
	log.Println(e.Error)
	return json.Marshal(struct {
		Error string `json:"error"`
	}{e.Error.Error()})
}

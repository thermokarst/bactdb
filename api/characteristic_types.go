package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/thermokarst/bactdb/models"
)

func serveCharacteristicType(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.ParseInt(mux.Vars(r)["Id"], 10, 0)
	if err != nil {
		return err
	}

	characteristic_type, err := store.CharacteristicTypes.Get(id)
	if err != nil {
		return err
	}

	return writeJSON(w, characteristic_type)
}

func serveCreateCharacteristicType(w http.ResponseWriter, r *http.Request) error {
	var characteristic_type models.CharacteristicType
	err := json.NewDecoder(r.Body).Decode(&characteristic_type)
	if err != nil {
		return err
	}

	created, err := store.CharacteristicTypes.Create(&characteristic_type)
	if err != nil {
		return err
	}
	if created {
		w.WriteHeader(http.StatusCreated)
	}

	return writeJSON(w, characteristic_type)
}

func serveCharacteristicTypeList(w http.ResponseWriter, r *http.Request) error {
	var opt models.CharacteristicTypeListOptions
	if err := schemaDecoder.Decode(&opt, r.URL.Query()); err != nil {
		return err
	}

	characteristic_types, err := store.CharacteristicTypes.List(&opt)
	if err != nil {
		return err
	}
	if characteristic_types == nil {
		characteristic_types = []*models.CharacteristicType{}
	}

	return writeJSON(w, characteristic_types)
}

func serveUpdateCharacteristicType(w http.ResponseWriter, r *http.Request) error {
	id, _ := strconv.ParseInt(mux.Vars(r)["Id"], 10, 0)
	var characteristic_type models.CharacteristicType
	err := json.NewDecoder(r.Body).Decode(&characteristic_type)
	if err != nil {
		return err
	}

	updated, err := store.CharacteristicTypes.Update(id, &characteristic_type)
	if err != nil {
		return err
	}
	if updated {
		w.WriteHeader(http.StatusOK)
	}

	return writeJSON(w, characteristic_type)
}

func serveDeleteCharacteristicType(w http.ResponseWriter, r *http.Request) error {
	id, _ := strconv.ParseInt(mux.Vars(r)["Id"], 10, 0)

	deleted, err := store.CharacteristicTypes.Delete(id)
	if err != nil {
		return err
	}
	if deleted {
		w.WriteHeader(http.StatusOK)
	}

	return writeJSON(w, &models.CharacteristicType{})
}

package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/thermokarst/bactdb/models"
)

func serveCharacteristic(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.ParseInt(mux.Vars(r)["Id"], 10, 0)
	if err != nil {
		return err
	}

	characteristic, err := store.Characteristics.Get(id)
	if err != nil {
		return err
	}

	return writeJSON(w, characteristic)
}

func serveCreateCharacteristic(w http.ResponseWriter, r *http.Request) error {
	var characteristic models.Characteristic
	err := json.NewDecoder(r.Body).Decode(&characteristic)
	if err != nil {
		return err
	}

	created, err := store.Characteristics.Create(&characteristic)
	if err != nil {
		return err
	}
	if created {
		w.WriteHeader(http.StatusCreated)
	}

	return writeJSON(w, characteristic)
}

func serveCharacteristicList(w http.ResponseWriter, r *http.Request) error {
	var opt models.CharacteristicListOptions
	if err := schemaDecoder.Decode(&opt, r.URL.Query()); err != nil {
		return err
	}

	characteristics, err := store.Characteristics.List(&opt)
	if err != nil {
		return err
	}
	if characteristics == nil {
		characteristics = []*models.Characteristic{}
	}

	return writeJSON(w, characteristics)
}

func serveUpdateCharacteristic(w http.ResponseWriter, r *http.Request) error {
	id, _ := strconv.ParseInt(mux.Vars(r)["Id"], 10, 0)
	var characteristic models.Characteristic
	err := json.NewDecoder(r.Body).Decode(&characteristic)
	if err != nil {
		return err
	}

	updated, err := store.Characteristics.Update(id, &characteristic)
	if err != nil {
		return err
	}
	if updated {
		w.WriteHeader(http.StatusOK)
	}

	return writeJSON(w, characteristic)
}

func serveDeleteCharacteristic(w http.ResponseWriter, r *http.Request) error {
	id, _ := strconv.ParseInt(mux.Vars(r)["Id"], 10, 0)

	deleted, err := store.Characteristics.Delete(id)
	if err != nil {
		return err
	}
	if deleted {
		w.WriteHeader(http.StatusOK)
	}

	return writeJSON(w, &models.Characteristic{})
}

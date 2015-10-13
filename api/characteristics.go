package api

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/thermokarst/bactdb/errors"
	"github.com/thermokarst/bactdb/helpers"
	"github.com/thermokarst/bactdb/models"
	"github.com/thermokarst/bactdb/payloads"
	"github.com/thermokarst/bactdb/types"
)

// CharacteristicService provides for CRUD operations
type CharacteristicService struct{}

// Unmarshal satisfies interface Updater and interface Creater
func (c CharacteristicService) Unmarshal(b []byte) (types.Entity, error) {
	var cj payloads.Characteristic
	err := json.Unmarshal(b, &cj)
	return &cj, err
}

// List lists all characteristics
func (c CharacteristicService) List(val *url.Values, claims *types.Claims) (types.Entity, *types.AppError) {
	if val == nil {
		return nil, newJSONError(errors.ErrMustProvideOptions, http.StatusInternalServerError)
	}
	var opt helpers.ListOptions
	if err := helpers.SchemaDecoder.Decode(&opt, *val); err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	characteristics, err := models.ListCharacteristics(opt, claims)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	strainsOpt, err := models.StrainOptsFromCharacteristics(opt)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	strains, err := models.ListStrains(*strainsOpt, claims)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	speciesOpt, err := models.SpeciesOptsFromStrains(*strainsOpt)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	species, err := models.ListSpecies(*speciesOpt, claims)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	measurementsOpt, err := models.MeasurementOptsFromCharacteristics(opt)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	measurements, err := models.ListMeasurements(*measurementsOpt, claims)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	payload := payloads.Characteristics{
		Characteristics: characteristics,
		Measurements:    measurements,
		Strains:         strains,
		Species:         species,
		Meta: &models.CharacteristicMeta{
			CanAdd: helpers.CanAdd(claims),
		},
	}

	return &payload, nil
}

// Get retrieves a single characteristic
func (c CharacteristicService) Get(id int64, genus string, claims *types.Claims) (types.Entity, *types.AppError) {
	characteristic, err := models.GetCharacteristic(id, genus, claims)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	strains, strainOpts, err := models.StrainsFromCharacteristicID(id, genus, claims)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	speciesOpt, err := models.SpeciesOptsFromStrains(*strainOpts)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	species, err := models.ListSpecies(*speciesOpt, claims)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	measurements, _, err := models.MeasurementsFromCharacteristicID(id, genus, claims)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	payload := payloads.Characteristic{
		Characteristic: characteristic,
		Measurements:   measurements,
		Strains:        strains,
		Species:        species,
	}

	return &payload, nil
}

// Update modifies an existing characteristic
func (c CharacteristicService) Update(id int64, e *types.Entity, genus string, claims *types.Claims) *types.AppError {
	payload := (*e).(*payloads.Characteristic)
	payload.Characteristic.UpdatedBy = claims.Sub
	payload.Characteristic.ID = id

	// First, handle Characteristic Type
	id, err := models.InsertOrGetCharacteristicType(payload.Characteristic.CharacteristicType, claims)
	if err != nil {
		return newJSONError(err, http.StatusInternalServerError)
	}

	payload.Characteristic.CanEdit = helpers.CanEdit(claims, payload.Characteristic.CreatedBy)

	payload.Characteristic.CharacteristicTypeID = id

	if err := models.Update(payload.Characteristic.CharacteristicBase); err != nil {
		if err == errors.ErrCharacteristicNotUpdated {
			return newJSONError(err, http.StatusBadRequest)
		}
		return newJSONError(err, http.StatusInternalServerError)
	}

	strains, strainOpts, err := models.StrainsFromCharacteristicID(id, genus, claims)
	if err != nil {
		return newJSONError(err, http.StatusInternalServerError)
	}

	speciesOpt, err := models.SpeciesOptsFromStrains(*strainOpts)
	if err != nil {
		return newJSONError(err, http.StatusInternalServerError)
	}

	species, err := models.ListSpecies(*speciesOpt, claims)
	if err != nil {
		return newJSONError(err, http.StatusInternalServerError)
	}

	payload.Strains = strains
	// TODO: tack on measurements
	payload.Measurements = nil
	payload.Species = species

	return nil
}

// Create initializes a new characteristic
func (c CharacteristicService) Create(e *types.Entity, genus string, claims *types.Claims) *types.AppError {
	payload := (*e).(*payloads.Characteristic)
	payload.Characteristic.CreatedBy = claims.Sub
	payload.Characteristic.UpdatedBy = claims.Sub

	id, err := models.InsertOrGetCharacteristicType(payload.Characteristic.CharacteristicType, claims)
	if err != nil {
		return newJSONError(err, http.StatusInternalServerError)
	}
	payload.Characteristic.CharacteristicTypeID = id

	if err := models.Create(payload.Characteristic.CharacteristicBase); err != nil {
		return newJSONError(err, http.StatusInternalServerError)
	}

	characteristic, err := models.GetCharacteristic(payload.Characteristic.ID, genus, claims)
	if err != nil {
		return newJSONError(err, http.StatusInternalServerError)
	}

	payload.Characteristic = characteristic
	payload.Meta = &models.CharacteristicMeta{
		CanAdd: helpers.CanAdd(claims),
	}
	return nil
}

// Delete deletes a single characteristic
func (c CharacteristicService) Delete(id int64, genus string, claims *types.Claims) *types.AppError {
	q := `DELETE FROM characteristics WHERE id=$1;`
	// TODO: fix this
	if _, err := models.DBH.Exec(q, id); err != nil {
		return newJSONError(err, http.StatusInternalServerError)
	}
	return nil
}

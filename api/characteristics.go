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

type CharacteristicService struct{}

func (c CharacteristicService) Unmarshal(b []byte) (types.Entity, error) {
	var cj payloads.Characteristic
	err := json.Unmarshal(b, &cj)
	return &cj, err
}

func (c CharacteristicService) List(val *url.Values, claims *types.Claims) (types.Entity, *types.AppError) {
	if val == nil {
		return nil, NewJSONError(errors.MustProvideOptions, http.StatusInternalServerError)
	}
	var opt helpers.ListOptions
	if err := helpers.SchemaDecoder.Decode(&opt, *val); err != nil {
		return nil, NewJSONError(err, http.StatusInternalServerError)
	}

	characteristics, err := models.ListCharacteristics(opt, claims)
	if err != nil {
		return nil, NewJSONError(err, http.StatusInternalServerError)
	}

	strains_opt, err := models.StrainOptsFromCharacteristics(opt)
	if err != nil {
		return nil, NewJSONError(err, http.StatusInternalServerError)
	}

	strains, err := models.ListStrains(*strains_opt, claims)
	if err != nil {
		return nil, NewJSONError(err, http.StatusInternalServerError)
	}

	species_opt, err := models.SpeciesOptsFromStrains(*strains_opt)
	if err != nil {
		return nil, NewJSONError(err, http.StatusInternalServerError)
	}

	species, err := models.ListSpecies(*species_opt, claims)
	if err != nil {
		return nil, NewJSONError(err, http.StatusInternalServerError)
	}

	measurements_opt, err := models.MeasurementOptsFromCharacteristics(opt)
	if err != nil {
		return nil, NewJSONError(err, http.StatusInternalServerError)
	}

	measurements, err := models.ListMeasurements(*measurements_opt, claims)
	if err != nil {
		return nil, NewJSONError(err, http.StatusInternalServerError)
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

func (c CharacteristicService) Get(id int64, genus string, claims *types.Claims) (types.Entity, *types.AppError) {
	characteristic, err := models.GetCharacteristic(id, genus, claims)
	if err != nil {
		return nil, NewJSONError(err, http.StatusInternalServerError)
	}

	strains, strain_opts, err := models.StrainsFromCharacteristicId(id, genus, claims)
	if err != nil {
		return nil, NewJSONError(err, http.StatusInternalServerError)
	}

	species_opt, err := models.SpeciesOptsFromStrains(*strain_opts)
	if err != nil {
		return nil, NewJSONError(err, http.StatusInternalServerError)
	}

	species, err := models.ListSpecies(*species_opt, claims)
	if err != nil {
		return nil, NewJSONError(err, http.StatusInternalServerError)
	}

	measurements, _, err := models.MeasurementsFromCharacteristicId(id, genus, claims)
	if err != nil {
		return nil, NewJSONError(err, http.StatusInternalServerError)
	}

	payload := payloads.Characteristic{
		Characteristic: characteristic,
		Measurements:   measurements,
		Strains:        strains,
		Species:        species,
	}

	return &payload, nil
}

func (c CharacteristicService) Update(id int64, e *types.Entity, genus string, claims *types.Claims) *types.AppError {
	payload := (*e).(*payloads.Characteristic)
	payload.Characteristic.UpdatedBy = claims.Sub
	payload.Characteristic.Id = id

	// First, handle Characteristic Type
	id, err := models.InsertOrGetCharacteristicType(payload.Characteristic.CharacteristicType, claims)
	if err != nil {
		return NewJSONError(err, http.StatusInternalServerError)
	}

	payload.Characteristic.CanEdit = helpers.CanEdit(claims, payload.Characteristic.CreatedBy)

	payload.Characteristic.CharacteristicTypeId = id
	// TODO: fix this
	count, err := models.DBH.Update(payload.Characteristic.CharacteristicBase)
	if err != nil {
		return NewJSONError(err, http.StatusInternalServerError)
	}
	if count != 1 {
		// TODO: fix this
		return NewJSONError(errors.CharacteristicNotUpdated, http.StatusBadRequest)
	}

	strains, strain_opts, err := models.StrainsFromCharacteristicId(id, genus, claims)
	if err != nil {
		return NewJSONError(err, http.StatusInternalServerError)
	}

	species_opt, err := models.SpeciesOptsFromStrains(*strain_opts)
	if err != nil {
		return NewJSONError(err, http.StatusInternalServerError)
	}

	species, err := models.ListSpecies(*species_opt, claims)
	if err != nil {
		return NewJSONError(err, http.StatusInternalServerError)
	}

	payload.Strains = strains
	// TODO: tack on measurements
	payload.Measurements = nil
	payload.Species = species

	return nil
}

func (c CharacteristicService) Create(e *types.Entity, genus string, claims *types.Claims) *types.AppError {
	payload := (*e).(*payloads.Characteristic)
	payload.Characteristic.CreatedBy = claims.Sub
	payload.Characteristic.UpdatedBy = claims.Sub

	id, err := models.InsertOrGetCharacteristicType(payload.Characteristic.CharacteristicType, claims)
	if err != nil {
		return NewJSONError(err, http.StatusInternalServerError)
	}
	payload.Characteristic.CharacteristicTypeId = id

	// TODO: fix this
	err = models.DBH.Insert(payload.Characteristic.CharacteristicBase)
	if err != nil {
		return NewJSONError(err, http.StatusInternalServerError)
	}

	characteristic, err := models.GetCharacteristic(payload.Characteristic.Id, genus, claims)
	if err != nil {
		return NewJSONError(err, http.StatusInternalServerError)
	}

	payload.Characteristic = characteristic
	payload.Meta = &models.CharacteristicMeta{
		CanAdd: helpers.CanAdd(claims),
	}
	return nil
}

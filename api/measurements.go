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

// MeasurementService provides for CRUD operations.
type MeasurementService struct{}

// Unmarshal satisfies interface Updater and interface Creater.
func (m MeasurementService) Unmarshal(b []byte) (types.Entity, error) {
	var mj payloads.Measurement
	err := json.Unmarshal(b, &mj)
	return &mj, err
}

// List lists all measurements.
func (m MeasurementService) List(val *url.Values, claims *types.Claims) (types.Entity, *types.AppError) {
	if val == nil {
		return nil, newJSONError(errors.ErrMustProvideOptions, http.StatusInternalServerError)
	}
	var opt helpers.MeasurementListOptions
	if err := helpers.SchemaDecoder.Decode(&opt, *val); err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	measurements, err := models.ListMeasurements(opt, claims)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	charOpts, err := models.CharacteristicOptsFromMeasurements(opt)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	characteristics, err := models.ListCharacteristics(*charOpts, claims)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	strainOpts, err := models.StrainOptsFromMeasurements(opt)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	strains, err := models.ListStrains(*strainOpts, claims)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	payload := payloads.Measurements{
		Characteristics: characteristics,
		Strains:         strains,
		Measurements:    measurements,
	}

	return &payload, nil
}

// Get retrieves a single measurement.
func (m MeasurementService) Get(id int64, genus string, claims *types.Claims) (types.Entity, *types.AppError) {
	measurement, err := models.GetMeasurement(id, genus, claims)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	payload := payloads.Measurement{
		Measurement: measurement,
	}

	return &payload, nil
}

// Update modifies a single measurement.
func (m MeasurementService) Update(id int64, e *types.Entity, genus string, claims *types.Claims) *types.AppError {
	payload := (*e).(*payloads.Measurement)
	payload.Measurement.UpdatedBy = claims.Sub
	payload.Measurement.ID = id

	if payload.Measurement.TextMeasurementType.Valid {
		id, err := models.GetTextMeasurementTypeID(payload.Measurement.TextMeasurementType.String)
		if err != nil {
			return newJSONError(err, http.StatusInternalServerError)
		}
		payload.Measurement.TextMeasurementTypeID.Int64 = id
		payload.Measurement.TextMeasurementTypeID.Valid = true
	}

	if err := models.Update(payload.Measurement.MeasurementBase); err != nil {
		if err == errors.ErrMeasurementNotUpdated {
			return newJSONError(err, http.StatusBadRequest)
		}
		return newJSONError(err, http.StatusInternalServerError)
	}

	measurement, err := models.GetMeasurement(id, genus, claims)
	if err != nil {
		return newJSONError(err, http.StatusInternalServerError)
	}

	payload.Measurement = measurement

	return nil
}

// Delete deletes a single measurement.
func (m MeasurementService) Delete(id int64, genus string, claims *types.Claims) *types.AppError {
	measurement, err := models.GetMeasurement(id, genus, claims)
	if err != nil {
		return newJSONError(err, http.StatusInternalServerError)
	}
	if err := models.Delete(measurement); err != nil {
		return newJSONError(err, http.StatusInternalServerError)
	}

	return nil
}

// Create initializes a new measurement.
func (m MeasurementService) Create(e *types.Entity, genus string, claims *types.Claims) *types.AppError {
	payload := (*e).(*payloads.Measurement)
	payload.Measurement.CreatedBy = claims.Sub
	payload.Measurement.UpdatedBy = claims.Sub

	if err := models.Create(payload.Measurement.MeasurementBase); err != nil {
		return newJSONError(err, http.StatusInternalServerError)
	}

	return nil

}

package api

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/thermokarst/bactdb/helpers"
	"github.com/thermokarst/bactdb/models"
	"github.com/thermokarst/bactdb/payloads"
	"github.com/thermokarst/bactdb/types"
)

type MeasurementService struct{}

func (s MeasurementService) Unmarshal(b []byte) (types.Entity, error) {
	var mj payloads.MeasurementPayload
	err := json.Unmarshal(b, &mj)
	return &mj, err
}

func (m MeasurementService) List(val *url.Values, claims *types.Claims) (types.Entity, *types.AppError) {
	if val == nil {
		return nil, helpers.ErrMustProvideOptionsJSON
	}
	var opt helpers.MeasurementListOptions
	if err := helpers.SchemaDecoder.Decode(&opt, *val); err != nil {
		return nil, types.NewJSONError(err, http.StatusInternalServerError)
	}

	measurements, err := models.ListMeasurements(opt, claims)
	if err != nil {
		return nil, types.NewJSONError(err, http.StatusInternalServerError)
	}

	char_opts, err := models.CharacteristicOptsFromMeasurements(opt)
	if err != nil {
		return nil, types.NewJSONError(err, http.StatusInternalServerError)
	}

	characteristics, err := models.ListCharacteristics(*char_opts, claims)
	if err != nil {
		return nil, types.NewJSONError(err, http.StatusInternalServerError)
	}

	strain_opts, err := models.StrainOptsFromMeasurements(opt)
	if err != nil {
		return nil, types.NewJSONError(err, http.StatusInternalServerError)
	}

	strains, err := models.ListStrains(*strain_opts, claims)
	if err != nil {
		return nil, types.NewJSONError(err, http.StatusInternalServerError)
	}

	payload := payloads.MeasurementsPayload{
		Characteristics: characteristics,
		Strains:         strains,
		Measurements:    measurements,
	}

	return &payload, nil
}

func (m MeasurementService) Get(id int64, genus string, claims *types.Claims) (types.Entity, *types.AppError) {
	measurement, err := models.GetMeasurement(id, genus, claims)
	if err != nil {
		return nil, types.NewJSONError(err, http.StatusInternalServerError)
	}

	payload := payloads.MeasurementPayload{
		Measurement: measurement,
	}

	return &payload, nil
}

func (s MeasurementService) Update(id int64, e *types.Entity, genus string, claims *types.Claims) *types.AppError {
	payload := (*e).(*payloads.MeasurementPayload)
	payload.Measurement.UpdatedBy = claims.Sub
	payload.Measurement.Id = id

	if payload.Measurement.TextMeasurementType.Valid {
		id, err := models.GetTextMeasurementTypeId(payload.Measurement.TextMeasurementType.String)
		if err != nil {
			return types.NewJSONError(err, http.StatusInternalServerError)
		}
		payload.Measurement.TextMeasurementTypeId.Int64 = id
		payload.Measurement.TextMeasurementTypeId.Valid = true
	}

	// TODO: fix this
	count, err := models.DBH.Update(payload.Measurement.MeasurementBase)
	if err != nil {
		return types.NewJSONError(err, http.StatusInternalServerError)
	}
	if count != 1 {
		// TODO: fix this
		return types.NewJSONError(models.ErrStrainNotUpdated, http.StatusBadRequest)
	}

	measurement, err := models.GetMeasurement(id, genus, claims)
	if err != nil {
		return types.NewJSONError(err, http.StatusInternalServerError)
	}

	payload.Measurement = measurement

	return nil
}

func (m MeasurementService) Delete(id int64, genus string, claims *types.Claims) *types.AppError {
	q := `DELETE FROM measurements WHERE id=$1;`
	// TODO: fix this
	_, err := models.DBH.Exec(q, id)
	if err != nil {
		return types.NewJSONError(err, http.StatusInternalServerError)
	}
	return nil
}

func (m MeasurementService) Create(e *types.Entity, genus string, claims *types.Claims) *types.AppError {
	payload := (*e).(*payloads.MeasurementPayload)
	payload.Measurement.CreatedBy = claims.Sub
	payload.Measurement.UpdatedBy = claims.Sub

	// TODO: fix this
	if err := models.DBH.Insert(payload.Measurement.MeasurementBase); err != nil {
		return types.NewJSONError(err, http.StatusInternalServerError)
	}

	return nil

}

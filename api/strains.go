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

// StrainService provides for CRUD operations
type StrainService struct{}

// Unmarshal satisfies interface Updater and interface Creater
func (s StrainService) Unmarshal(b []byte) (types.Entity, error) {
	var sj payloads.Strain
	err := json.Unmarshal(b, &sj)
	return &sj, err
}

// List lists all strains
func (s StrainService) List(val *url.Values, claims *types.Claims) (types.Entity, *types.AppError) {
	if val == nil {
		return nil, newJSONError(errors.ErrMustProvideOptions, http.StatusInternalServerError)
	}
	var opt helpers.ListOptions
	if err := helpers.SchemaDecoder.Decode(&opt, *val); err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	strains, err := models.ListStrains(opt, claims)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	speciesOpt, err := models.SpeciesOptsFromStrains(opt)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	species, err := models.ListSpecies(*speciesOpt, claims)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	characteristicsOpt, err := models.CharacteristicsOptsFromStrains(opt)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	characteristics, err := models.ListCharacteristics(*characteristicsOpt, claims)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	characteristicIDs := []int64{}
	for _, c := range *characteristics {
		characteristicIDs = append(characteristicIDs, c.ID)
	}

	strainIDs := []int64{}
	for _, s := range *strains {
		strainIDs = append(strainIDs, s.ID)
	}

	measurementOpt := helpers.MeasurementListOptions{
		ListOptions: helpers.ListOptions{
			Genus: opt.Genus,
		},
		Strains:         strainIDs,
		Characteristics: characteristicIDs,
	}

	measurements, err := models.ListMeasurements(measurementOpt, claims)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	payload := payloads.Strains{
		Strains:         strains,
		Species:         species,
		Measurements:    measurements,
		Characteristics: characteristics,
		Meta: &models.StrainMeta{
			CanAdd: helpers.CanAdd(claims),
		},
	}

	return &payload, nil
}

// Get retrieves a single strain
func (s StrainService) Get(id int64, genus string, claims *types.Claims) (types.Entity, *types.AppError) {
	strain, err := models.GetStrain(id, genus, claims)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	species, err := models.GetSpecies(strain.SpeciesID, genus, claims)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	opt := helpers.ListOptions{Genus: genus, IDs: []int64{id}}
	characteristicsOpt, err := models.CharacteristicsOptsFromStrains(opt)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	characteristics, err := models.ListCharacteristics(*characteristicsOpt, claims)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	characteristicIDs := []int64{}
	for _, c := range *characteristics {
		characteristicIDs = append(characteristicIDs, c.ID)
	}

	measurementOpt := helpers.MeasurementListOptions{
		ListOptions: helpers.ListOptions{
			Genus: genus,
		},
		Strains:         []int64{id},
		Characteristics: characteristicIDs,
	}

	measurements, err := models.ListMeasurements(measurementOpt, claims)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	var manySpecies models.ManySpecies = []*models.Species{species}

	payload := payloads.Strain{
		Strain:          strain,
		Species:         &manySpecies,
		Characteristics: characteristics,
		Measurements:    measurements,
		Meta: &models.StrainMeta{
			CanAdd: helpers.CanAdd(claims),
		},
	}

	return &payload, nil
}

// Update modifies an existing strain
func (s StrainService) Update(id int64, e *types.Entity, genus string, claims *types.Claims) *types.AppError {
	payload := (*e).(*payloads.Strain)
	payload.Strain.UpdatedBy = claims.Sub
	payload.Strain.ID = id

	// TODO: fix this
	count, err := models.DBH.Update(payload.Strain.StrainBase)
	if err != nil {
		return newJSONError(err, http.StatusInternalServerError)
	}
	if count != 1 {
		// TODO: fix this
		return newJSONError(errors.ErrStrainNotUpdated, http.StatusBadRequest)
	}

	strain, err := models.GetStrain(id, genus, claims)
	if err != nil {
		return newJSONError(err, http.StatusInternalServerError)
	}

	species, err := models.GetSpecies(strain.SpeciesID, genus, claims)
	if err != nil {
		return newJSONError(err, http.StatusInternalServerError)
	}

	var manySpecies models.ManySpecies = []*models.Species{species}

	payload.Strain = strain
	payload.Species = &manySpecies
	payload.Meta = &models.StrainMeta{
		CanAdd: helpers.CanAdd(claims),
	}

	return nil
}

// Create initializes a new strain
func (s StrainService) Create(e *types.Entity, genus string, claims *types.Claims) *types.AppError {
	payload := (*e).(*payloads.Strain)
	payload.Strain.CreatedBy = claims.Sub
	payload.Strain.UpdatedBy = claims.Sub

	// TODO: fix this
	if err := models.DBH.Insert(payload.Strain.StrainBase); err != nil {
		return newJSONError(err, http.StatusInternalServerError)
	}

	strain, err := models.GetStrain(payload.Strain.ID, genus, claims)
	if err != nil {
		return newJSONError(err, http.StatusInternalServerError)
	}

	species, err := models.GetSpecies(strain.SpeciesID, genus, claims)
	if err != nil {
		return newJSONError(err, http.StatusInternalServerError)
	}

	var manySpecies models.ManySpecies = []*models.Species{species}

	payload.Strain = strain
	payload.Species = &manySpecies
	payload.Meta = &models.StrainMeta{
		CanAdd: helpers.CanAdd(claims),
	}

	return nil
}

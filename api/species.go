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

// SpeciesService provides for CRUD operations
type SpeciesService struct{}

// Unmarshal satisfies interface Updater and interface Creater
func (s SpeciesService) Unmarshal(b []byte) (types.Entity, error) {
	var sj payloads.Species
	err := json.Unmarshal(b, &sj)
	return &sj, err
}

// List lists species
func (s SpeciesService) List(val *url.Values, claims *types.Claims) (types.Entity, *types.AppError) {
	if val == nil {
		return nil, newJSONError(errors.ErrMustProvideOptions, http.StatusInternalServerError)
	}
	var opt helpers.ListOptions
	if err := helpers.SchemaDecoder.Decode(&opt, *val); err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	species, err := models.ListSpecies(opt, claims)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	strainsOpt, err := models.StrainOptsFromSpecies(opt)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	strains, err := models.ListStrains(*strainsOpt, claims)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	payload := payloads.ManySpecies{
		Species: species,
		Strains: strains,
		Meta: &models.SpeciesMeta{
			CanAdd: helpers.CanAdd(claims),
		},
	}

	return &payload, nil
}

// Get retrieves a single species
func (s SpeciesService) Get(id int64, genus string, claims *types.Claims) (types.Entity, *types.AppError) {
	species, err := models.GetSpecies(id, genus, claims)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	strains, err := models.StrainsFromSpeciesID(id, genus, claims)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	payload := payloads.Species{
		Species: species,
		Strains: strains,
		Meta: &models.SpeciesMeta{
			CanAdd: helpers.CanAdd(claims),
		},
	}

	return &payload, nil
}

// Update modifies an existing species
func (s SpeciesService) Update(id int64, e *types.Entity, genus string, claims *types.Claims) *types.AppError {
	payload := (*e).(*payloads.Species)
	payload.Species.UpdatedBy = claims.Sub
	payload.Species.ID = id

	genusID, err := models.GenusIDFromName(genus)
	if err != nil {
		return newJSONError(err, http.StatusInternalServerError)
	}
	payload.Species.SpeciesBase.GenusID = genusID

	if err := models.Update(payload.Species.SpeciesBase); err != nil {
		if err == errors.ErrSpeciesNotUpdated {
			return newJSONError(err, http.StatusBadRequest)
		}
		return newJSONError(err, http.StatusInternalServerError)
	}

	// Reload to send back down the wire
	species, err := models.GetSpecies(id, genus, claims)
	if err != nil {
		return newJSONError(err, http.StatusInternalServerError)
	}

	strains, err := models.StrainsFromSpeciesID(id, genus, claims)
	if err != nil {
		return newJSONError(err, http.StatusInternalServerError)
	}

	payload.Species = species
	payload.Strains = strains
	payload.Meta = &models.SpeciesMeta{
		CanAdd: helpers.CanAdd(claims),
	}

	return nil
}

// Create initializes a new species
func (s SpeciesService) Create(e *types.Entity, genus string, claims *types.Claims) *types.AppError {
	payload := (*e).(*payloads.Species)
	payload.Species.CreatedBy = claims.Sub
	payload.Species.UpdatedBy = claims.Sub

	genusID, err := models.GenusIDFromName(genus)
	if err != nil {
		return newJSONError(err, http.StatusInternalServerError)
	}
	payload.Species.SpeciesBase.GenusID = genusID

	if err := models.Create(payload.Species.SpeciesBase); err != nil {
		return newJSONError(err, http.StatusInternalServerError)
	}

	// Reload to send back down the wire
	species, err := models.GetSpecies(payload.Species.ID, genus, claims)
	if err != nil {
		return newJSONError(err, http.StatusInternalServerError)
	}

	// Note, no strains when new species

	payload.Species = species
	payload.Meta = &models.SpeciesMeta{
		CanAdd: helpers.CanAdd(claims),
	}
	return nil
}

// Delete deletes a single species
func (s SpeciesService) Delete(id int64, genus string, claims *types.Claims) *types.AppError {
	q := `DELETE FROM species WHERE id=$1;`
	// TODO: fix this
	if _, err := models.DBH.Exec(q, id); err != nil {
		return newJSONError(err, http.StatusInternalServerError)
	}
	return nil
}

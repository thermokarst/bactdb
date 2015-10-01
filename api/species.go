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

type SpeciesService struct{}

func (s SpeciesService) Unmarshal(b []byte) (types.Entity, error) {
	var sj payloads.Species
	err := json.Unmarshal(b, &sj)
	return &sj, err
}

func (s SpeciesService) List(val *url.Values, claims *types.Claims) (types.Entity, *types.AppError) {
	if val == nil {
		return nil, NewJSONError(errors.MustProvideOptions, http.StatusInternalServerError)
	}
	var opt helpers.ListOptions
	if err := helpers.SchemaDecoder.Decode(&opt, *val); err != nil {
		return nil, NewJSONError(err, http.StatusInternalServerError)
	}

	species, err := models.ListSpecies(opt, claims)
	if err != nil {
		return nil, NewJSONError(err, http.StatusInternalServerError)
	}

	strains_opt, err := models.StrainOptsFromSpecies(opt)
	if err != nil {
		return nil, NewJSONError(err, http.StatusInternalServerError)
	}

	strains, err := models.ListStrains(*strains_opt, claims)
	if err != nil {
		return nil, NewJSONError(err, http.StatusInternalServerError)
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

func (s SpeciesService) Get(id int64, genus string, claims *types.Claims) (types.Entity, *types.AppError) {
	species, err := models.GetSpecies(id, genus, claims)
	if err != nil {
		return nil, NewJSONError(err, http.StatusInternalServerError)
	}

	strains, err := models.StrainsFromSpeciesId(id, genus, claims)
	if err != nil {
		return nil, NewJSONError(err, http.StatusInternalServerError)
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

func (s SpeciesService) Update(id int64, e *types.Entity, genus string, claims *types.Claims) *types.AppError {
	payload := (*e).(*payloads.Species)
	payload.Species.UpdatedBy = claims.Sub
	payload.Species.Id = id

	genus_id, err := models.GenusIdFromName(genus)
	if err != nil {
		return NewJSONError(err, http.StatusInternalServerError)
	}
	payload.Species.SpeciesBase.GenusID = genus_id

	// TODO: fix this
	count, err := models.DBH.Update(payload.Species.SpeciesBase)
	if err != nil {
		return NewJSONError(err, http.StatusInternalServerError)
	}
	if count != 1 {
		// TODO: fix this
		return NewJSONError(errors.SpeciesNotUpdated, http.StatusBadRequest)
	}

	// Reload to send back down the wire
	species, err := models.GetSpecies(id, genus, claims)
	if err != nil {
		return NewJSONError(err, http.StatusInternalServerError)
	}

	strains, err := models.StrainsFromSpeciesId(id, genus, claims)
	if err != nil {
		return NewJSONError(err, http.StatusInternalServerError)
	}

	payload.Species = species
	payload.Strains = strains
	payload.Meta = &models.SpeciesMeta{
		CanAdd: helpers.CanAdd(claims),
	}

	return nil
}

func (s SpeciesService) Create(e *types.Entity, genus string, claims *types.Claims) *types.AppError {
	payload := (*e).(*payloads.Species)
	payload.Species.CreatedBy = claims.Sub
	payload.Species.UpdatedBy = claims.Sub

	genus_id, err := models.GenusIdFromName(genus)
	if err != nil {
		return NewJSONError(err, http.StatusInternalServerError)
	}
	payload.Species.SpeciesBase.GenusID = genus_id

	// TODO: fix this
	err = models.DBH.Insert(payload.Species.SpeciesBase)
	if err != nil {
		return NewJSONError(err, http.StatusInternalServerError)
	}

	// Reload to send back down the wire
	species, err := models.GetSpecies(payload.Species.Id, genus, claims)
	if err != nil {
		return NewJSONError(err, http.StatusInternalServerError)
	}

	// Note, no strains when new species

	payload.Species = species
	payload.Meta = &models.SpeciesMeta{
		CanAdd: helpers.CanAdd(claims),
	}
	return nil
}

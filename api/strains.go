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

type StrainService struct{}

func (s StrainService) Unmarshal(b []byte) (types.Entity, error) {
	var sj payloads.Strain
	err := json.Unmarshal(b, &sj)
	return &sj, err
}

func (s StrainService) List(val *url.Values, claims *types.Claims) (types.Entity, *types.AppError) {
	if val == nil {
		return nil, helpers.ErrMustProvideOptionsJSON
	}
	var opt helpers.ListOptions
	if err := helpers.SchemaDecoder.Decode(&opt, *val); err != nil {
		return nil, types.NewJSONError(err, http.StatusInternalServerError)
	}

	strains, err := models.ListStrains(opt, claims)
	if err != nil {
		return nil, types.NewJSONError(err, http.StatusInternalServerError)
	}

	species_opt, err := models.SpeciesOptsFromStrains(opt)
	if err != nil {
		return nil, types.NewJSONError(err, http.StatusInternalServerError)
	}

	species, err := models.ListSpecies(*species_opt, claims)
	if err != nil {
		return nil, types.NewJSONError(err, http.StatusInternalServerError)
	}

	characteristics_opt, err := models.CharacteristicsOptsFromStrains(opt)
	if err != nil {
		return nil, types.NewJSONError(err, http.StatusInternalServerError)
	}

	characteristics, err := models.ListCharacteristics(*characteristics_opt, claims)
	if err != nil {
		return nil, types.NewJSONError(err, http.StatusInternalServerError)
	}

	characteristic_ids := []int64{}
	for _, c := range *characteristics {
		characteristic_ids = append(characteristic_ids, c.Id)
	}

	strain_ids := []int64{}
	for _, s := range *strains {
		strain_ids = append(strain_ids, s.Id)
	}

	measurement_opt := helpers.MeasurementListOptions{
		ListOptions: helpers.ListOptions{
			Genus: opt.Genus,
		},
		Strains:         strain_ids,
		Characteristics: characteristic_ids,
	}

	measurements, err := models.ListMeasurements(measurement_opt, claims)
	if err != nil {
		return nil, types.NewJSONError(err, http.StatusInternalServerError)
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

func (s StrainService) Get(id int64, genus string, claims *types.Claims) (types.Entity, *types.AppError) {
	strain, err := models.GetStrain(id, genus, claims)
	if err != nil {
		return nil, types.NewJSONError(err, http.StatusInternalServerError)
	}

	species, err := models.GetSpecies(strain.SpeciesId, genus, claims)
	if err != nil {
		return nil, types.NewJSONError(err, http.StatusInternalServerError)
	}

	opt := helpers.ListOptions{Genus: genus, Ids: []int64{id}}
	characteristics_opt, err := models.CharacteristicsOptsFromStrains(opt)
	if err != nil {
		return nil, types.NewJSONError(err, http.StatusInternalServerError)
	}

	characteristics, err := models.ListCharacteristics(*characteristics_opt, claims)
	if err != nil {
		return nil, types.NewJSONError(err, http.StatusInternalServerError)
	}

	characteristic_ids := []int64{}
	for _, c := range *characteristics {
		characteristic_ids = append(characteristic_ids, c.Id)
	}

	measurement_opt := helpers.MeasurementListOptions{
		ListOptions: helpers.ListOptions{
			Genus: genus,
		},
		Strains:         []int64{id},
		Characteristics: characteristic_ids,
	}

	measurements, err := models.ListMeasurements(measurement_opt, claims)
	if err != nil {
		return nil, types.NewJSONError(err, http.StatusInternalServerError)
	}

	var many_species models.ManySpecies = []*models.Species{species}

	payload := payloads.Strain{
		Strain:          strain,
		Species:         &many_species,
		Characteristics: characteristics,
		Measurements:    measurements,
		Meta: &models.StrainMeta{
			CanAdd: helpers.CanAdd(claims),
		},
	}

	return &payload, nil
}

func (s StrainService) Update(id int64, e *types.Entity, genus string, claims *types.Claims) *types.AppError {
	payload := (*e).(*payloads.Strain)
	payload.Strain.UpdatedBy = claims.Sub
	payload.Strain.Id = id

	// TODO: fix this
	count, err := models.DBH.Update(payload.Strain.StrainBase)
	if err != nil {
		return types.NewJSONError(err, http.StatusInternalServerError)
	}
	if count != 1 {
		// TODO: fix this
		return types.NewJSONError(models.ErrStrainNotUpdated, http.StatusBadRequest)
	}

	strain, err := models.GetStrain(id, genus, claims)
	if err != nil {
		return types.NewJSONError(err, http.StatusInternalServerError)
	}

	species, err := models.GetSpecies(strain.SpeciesId, genus, claims)
	if err != nil {
		return types.NewJSONError(err, http.StatusInternalServerError)
	}

	var many_species models.ManySpecies = []*models.Species{species}

	payload.Strain = strain
	payload.Species = &many_species
	payload.Meta = &models.StrainMeta{
		CanAdd: helpers.CanAdd(claims),
	}

	return nil
}

func (s StrainService) Create(e *types.Entity, genus string, claims *types.Claims) *types.AppError {
	payload := (*e).(*payloads.Strain)
	payload.Strain.CreatedBy = claims.Sub
	payload.Strain.UpdatedBy = claims.Sub

	// TODO: fix this
	if err := models.DBH.Insert(payload.Strain.StrainBase); err != nil {
		return types.NewJSONError(err, http.StatusInternalServerError)
	}

	strain, err := models.GetStrain(payload.Strain.Id, genus, claims)
	if err != nil {
		return types.NewJSONError(err, http.StatusInternalServerError)
	}

	species, err := models.GetSpecies(strain.SpeciesId, genus, claims)
	if err != nil {
		return types.NewJSONError(err, http.StatusInternalServerError)
	}

	var many_species models.ManySpecies = []*models.Species{species}

	payload.Strain = strain
	payload.Species = &many_species
	payload.Meta = &models.StrainMeta{
		CanAdd: helpers.CanAdd(claims),
	}

	return nil
}

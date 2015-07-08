package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

var (
	ErrStrainNotFound       = errors.New("Strain not found")
	ErrStrainNotFoundJSON   = newJSONError(ErrStrainNotFound, http.StatusNotFound)
	ErrStrainNotUpdated     = errors.New("Strain not updated")
	ErrStrainNotUpdatedJSON = newJSONError(ErrStrainNotUpdated, http.StatusBadRequest)
)

func init() {
	DB.AddTableWithName(StrainBase{}, "strains").SetKeys(true, "Id")
}

type StrainService struct{}

// StrainBase is what the DB expects to see for inserts/updates
type StrainBase struct {
	Id                  int64      `db:"id" json:"id"`
	SpeciesId           int64      `db:"species_id" json:"species"`
	StrainName          string     `db:"strain_name" json:"strainName"`
	TypeStrain          bool       `db:"type_strain" json:"typeStrain"`
	AccessionNumbers    string     `db:"accession_numbers" json:"accessionNumbers"`
	Genbank             NullString `db:"genbank" json:"genbank"`
	WholeGenomeSequence NullString `db:"whole_genome_sequence" json:"wholeGenomeSequence"`
	IsolatedFrom        NullString `db:"isolated_from" json:"isolatedFrom"`
	Notes               NullString `db:"notes" json:"notes"`
	CreatedAt           NullTime   `db:"created_at" json:"createdAt"`
	UpdatedAt           NullTime   `db:"updated_at" json:"updatedAt"`
	DeletedAt           NullTime   `db:"deleted_at" json:"deletedAt"`
	CreatedBy           int64      `db:"created_by" json:"createdBy"`
	UpdatedBy           int64      `db:"updated_by" json:"updatedBy"`
	DeletedBy           NullInt64  `db:"deleted_by" json:"deletedBy"`
}

type Strain struct {
	*StrainBase
	Measurements      NullSliceInt64 `db:"measurements" json:"measurements"`
	TotalMeasurements int64          `db:"total_measurements" json:"totalMeasurements"`
	SortOrder         int64          `db:"sort_order" json:"sortOrder"`
}

type Strains []*Strain

type StrainMeta struct {
	CanAdd  bool    `json:"canAdd"`
	CanEdit []int64 `json:"canEdit"`
}

type StrainPayload struct {
	Strain  *Strain      `json:"strain"`
	Species *ManySpecies `json:"species"`
	Meta    *StrainMeta  `json:"meta"`
}

type StrainsPayload struct {
	Strains *Strains     `json:"strains"`
	Species *ManySpecies `json:"species"`
	Meta    *StrainMeta  `json:"meta"`
}

func (s *StrainPayload) marshal() ([]byte, error) {
	return json.Marshal(s)
}

func (s *StrainsPayload) marshal() ([]byte, error) {
	return json.Marshal(s)
}

func (s StrainService) unmarshal(b []byte) (entity, error) {
	var sj StrainPayload
	err := json.Unmarshal(b, &sj)
	return &sj, err
}

func (s StrainService) list(val *url.Values, claims Claims) (entity, *appError) {
	if val == nil {
		return nil, ErrMustProvideOptionsJSON
	}
	var opt ListOptions
	if err := schemaDecoder.Decode(&opt, *val); err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	strains, err := listStrains(opt)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	species_opt, err := speciesOptsFromStrains(opt)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	species, err := listSpecies(*species_opt)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	edit_list := make(map[int64]int64)

	for _, v := range *strains {
		edit_list[v.Id] = v.CreatedBy
	}

	payload := StrainsPayload{
		Strains: strains,
		Species: species,
		Meta: &StrainMeta{
			CanAdd:  canAdd(claims),
			CanEdit: canEdit(claims, edit_list),
		},
	}

	return &payload, nil
}

func (s StrainService) get(id int64, genus string, claims Claims) (entity, *appError) {
	strain, err := getStrain(id, genus)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	species, err := getSpecies(strain.SpeciesId, genus)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	var many_species ManySpecies = []*Species{species}

	payload := StrainPayload{
		Strain:  strain,
		Species: &many_species,
		Meta: &StrainMeta{
			CanAdd:  canAdd(claims),
			CanEdit: canEdit(claims, map[int64]int64{strain.Id: strain.CreatedBy}),
		},
	}

	return &payload, nil
}

func (s StrainService) update(id int64, e *entity, claims Claims) *appError {
	payload := (*e).(*StrainPayload)
	payload.Strain.UpdatedBy = claims.Sub
	payload.Strain.UpdatedAt = currentTime()
	payload.Strain.Id = id

	count, err := DBH.Update(payload.Strain.StrainBase)
	if err != nil {
		return newJSONError(err, http.StatusInternalServerError)
	}
	if count != 1 {
		return ErrStrainNotUpdatedJSON
	}

	// TODO: add species and meta to payload

	return nil
}

func (s StrainService) create(e *entity, claims Claims) *appError {
	payload := (*e).(*StrainPayload)
	ct := currentTime()
	payload.Strain.CreatedBy = claims.Sub
	payload.Strain.CreatedAt = ct
	payload.Strain.UpdatedBy = claims.Sub
	payload.Strain.UpdatedAt = ct

	if err := DBH.Insert(payload.Strain.StrainBase); err != nil {
		return newJSONError(err, http.StatusInternalServerError)
	}

	// TODO: add species and meta to payload

	return nil
}

func listStrains(opt ListOptions) (*Strains, error) {
	var vals []interface{}

	q := `SELECT st.*, array_agg(m.id) AS measurements, COUNT(m) AS total_measurements,
		rank() OVER (ORDER BY sp.species_name ASC, st.type_strain ASC, st.strain_name ASC) AS sort_order
		FROM strains st
		INNER JOIN species sp ON sp.id=st.species_id
		INNER JOIN genera g ON g.id=sp.genus_id AND LOWER(g.genus_name)=LOWER($1)
		LEFT OUTER JOIN measurements m ON m.strain_id=st.id`
	vals = append(vals, opt.Genus)

	if len(opt.Ids) != 0 {
		var conds []string
		s := "st.id IN ("
		for i, id := range opt.Ids {
			s = s + fmt.Sprintf("$%v,", i+2) // start param index at 2
			vals = append(vals, id)
		}
		s = s[:len(s)-1] + ")"
		conds = append(conds, s)
		q += " WHERE (" + strings.Join(conds, ") AND (") + ")"
	}

	q += " GROUP BY st.id, st.species_id, sp.species_name;"

	strains := make(Strains, 0)
	err := DBH.Select(&strains, q, vals...)
	if err != nil {
		return nil, err
	}
	return &strains, nil
}

func getStrain(id int64, genus string) (*Strain, error) {
	var strain Strain
	q := `SELECT st.*, array_agg(m.id) AS measurements,
		COUNT(m) AS total_measurements, 0 AS sort_order
		FROM strains st
		INNER JOIN species sp ON sp.id=st.species_id
		INNER JOIN genera g ON g.id=sp.genus_id AND LOWER(g.genus_name)=LOWER($1)
		LEFT OUTER JOIN measurements m ON m.strain_id=st.id
		WHERE st.id=$2
		GROUP BY st.id, st.species_id;`
	if err := DBH.SelectOne(&strain, q, genus, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrStrainNotFound
		}
		return nil, err
	}
	return &strain, nil
}

func speciesOptsFromStrains(opt ListOptions) (*ListOptions, error) {
	relatedSpeciesIds := make([]int64, 0)

	if opt.Ids == nil {
		q := `SELECT DISTINCT st.species_id
			FROM strains st
			INNER JOIN species sp ON sp.id=st.species_id
			INNER JOIN genera g ON g.id=sp.genus_id AND LOWER(g.genus_name)=LOWER($1);`
		if err := DBH.Select(&relatedSpeciesIds, q, opt.Genus); err != nil {
			return nil, err
		}
	} else {
		var vals []interface{}
		var count int64 = 1
		q := fmt.Sprintf("SELECT DISTINCT species_id FROM strains WHERE %s;", valsIn("id", opt.Ids, &vals, &count))
		if err := DBH.Select(&relatedSpeciesIds, q, vals...); err != nil {
			return nil, err
		}
	}

	return &ListOptions{Genus: opt.Genus, Ids: relatedSpeciesIds}, nil
}

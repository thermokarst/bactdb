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
	ErrSpeciesNotFound       = errors.New("Species not found")
	ErrSpeciesNotFoundJSON   = newJSONError(ErrSpeciesNotFound, http.StatusNotFound)
	ErrSpeciesNotUpdated     = errors.New("Species not updated")
	ErrSpeciesNotUpdatedJSON = newJSONError(ErrSpeciesNotUpdated, http.StatusBadRequest)
)

func init() {
	DB.AddTableWithName(SpeciesBase{}, "species").SetKeys(true, "Id")
}

type SpeciesService struct{}

// SpeciesBase is what the DB expects to see for inserts/updates
type SpeciesBase struct {
	Id                  int64      `db:"id" json:"id"`
	GenusID             int64      `db:"genus_id" json:"-"`
	SubspeciesSpeciesID NullInt64  `db:"subspecies_species_id" json:"-"`
	SpeciesName         string     `db:"species_name" json:"speciesName"`
	TypeSpecies         NullBool   `db:"type_species" json:"typeSpecies"`
	Etymology           NullString `db:"etymology" json:"etymology"`
	CreatedAt           NullTime   `db:"created_at" json:"createdAt"`
	UpdatedAt           NullTime   `db:"updated_at" json:"updatedAt"`
	DeletedAt           NullTime   `db:"deleted_at" json:"deletedAt"`
	CreatedBy           int64      `db:"created_by" json:"createdBy"`
	UpdatedBy           int64      `db:"updated_by" json:"updatedBy"`
	DeletedBy           NullInt64  `db:"deleted_by" json:"deletedBy"`
}

// Species & SpeciesJSON(s) are what ember expects to see
type Species struct {
	*SpeciesBase
	GenusName    string         `db:"genus_name" json:"genusName"`
	Strains      NullSliceInt64 `db:"strains" json:"strains"`
	TotalStrains int64          `db:"total_strains" json:"totalStrains"`
	SortOrder    int64          `db:"sort_order" json:"sortOrder"`
}

type ManySpecies []*Species

type SpeciesMeta struct {
	CanAdd bool `json:"canAdd"`
}

type ManySpeciesPayload struct {
	Species *ManySpecies `json:"species"`
	Strains *Strains     `json:"strains"`
	Meta    *SpeciesMeta `json:"meta"`
}

type SpeciesJSON struct {
	Species *Species `json:"species"`
}

func (s *Species) marshal() ([]byte, error) {
	return json.Marshal(&SpeciesJSON{Species: s})
}

func (s *ManySpeciesPayload) marshal() ([]byte, error) {
	return json.Marshal(s)
}

func (s SpeciesService) unmarshal(b []byte) (entity, error) {
	var sj SpeciesJSON
	err := json.Unmarshal(b, &sj)
	return sj.Species, err
}

func (s SpeciesService) list(val *url.Values, claims Claims) (entity, *appError) {
	if val == nil {
		return nil, ErrMustProvideOptionsJSON
	}
	var opt ListOptions
	if err := schemaDecoder.Decode(&opt, *val); err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	var vals []interface{}
	sql := `SELECT sp.*, g.genus_name, array_agg(st.id) AS strains,
			COUNT(st) AS total_strains,
			rank() OVER (ORDER BY sp.species_name ASC) AS sort_order
			FROM species sp
			INNER JOIN genera g ON g.id=sp.genus_id AND LOWER(g.genus_name)=$1
			LEFT OUTER JOIN strains st ON st.species_id=sp.id`
	vals = append(vals, opt.Genus)

	if len(opt.Ids) != 0 {
		var conds []string
		s := "sp.id IN ("
		for i, id := range opt.Ids {
			s = s + fmt.Sprintf("$%v,", i+2) // start param index at 2
			vals = append(vals, id)
		}
		s = s[:len(s)-1] + ")"
		conds = append(conds, s)
		sql += " WHERE (" + strings.Join(conds, ") AND (") + ")"
	}

	sql += " GROUP BY sp.id, g.genus_name;"

	species := make(ManySpecies, 0)
	err := DBH.Select(&species, sql, vals...)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	strains_opt, err := strainOptsFromSpecies(opt)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	strains, err := listStrains(*strains_opt)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	payload := ManySpeciesPayload{
		Species: &species,
		Strains: strains,
		Meta:    &SpeciesMeta{CanAdd: canAdd(claims)},
	}

	return &payload, nil
}

func strainOptsFromSpecies(opt ListOptions) (*ListOptions, error) {
	relatedStrainIds := make([]int64, 0)

	if opt.Ids == nil {
		q := `SELECT DISTINCT st.id
			FROM strains st
			INNER JOIN species sp ON sp.id=st.species_id
			INNER JOIN genera g ON g.id=sp.genus_id AND LOWER(g.genus_name)=LOWER($1);`
		if err := DBH.Select(&relatedStrainIds, q, opt.Genus); err != nil {
			return nil, err
		}
	} else {
		var vals []interface{}
		var count int64 = 1
		q := fmt.Sprintf("SELECT DISTINCT id FROM strains WHERE %s;", valsIn("species_id", opt.Ids, &vals, &count))

		if err := DBH.Select(&relatedStrainIds, q, vals...); err != nil {
			return nil, err
		}
	}

	return &ListOptions{Genus: opt.Genus, Ids: relatedStrainIds}, nil
}

func (s SpeciesService) get(id int64, genus string) (entity, *appError) {
	var species Species
	q := `SELECT sp.*, g.genus_name, array_agg(st.id) AS strains,
		COUNT(st) AS total_strains, 0 AS sort_order
		FROM species sp
		INNER JOIN genera g ON g.id=sp.genus_id AND LOWER(g.genus_name)=$1
		LEFT OUTER JOIN strains st ON st.species_id=sp.id
		WHERE sp.id=$2
		GROUP BY sp.id, g.genus_name;`
	if err := DBH.SelectOne(&species, q, genus, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrSpeciesNotFoundJSON
		}
		return nil, newJSONError(err, http.StatusInternalServerError)
	}
	return &species, nil
}

func (s SpeciesService) update(id int64, e *entity, claims Claims) *appError {
	species := (*e).(*Species)
	species.UpdatedBy = claims.Sub
	species.UpdatedAt = currentTime()
	species.Id = id

	genus_id, err := genusIdFromName(species.GenusName)
	if err != nil {
		return newJSONError(err, http.StatusInternalServerError)
	}
	species.SpeciesBase.GenusID = genus_id

	count, err := DBH.Update(species.SpeciesBase)
	if err != nil {
		return newJSONError(err, http.StatusInternalServerError)
	}
	if count != 1 {
		return ErrSpeciesNotUpdatedJSON
	}
	return nil
}

func (s SpeciesService) create(e *entity, claims Claims) *appError {
	species := (*e).(*Species)
	ct := currentTime()
	species.CreatedBy = claims.Sub
	species.CreatedAt = ct
	species.UpdatedBy = claims.Sub
	species.UpdatedAt = ct

	genus_id, err := genusIdFromName(species.GenusName)
	if err != nil {
		return newJSONError(err, http.StatusInternalServerError)
	}
	species.SpeciesBase.GenusID = genus_id

	err = DBH.Insert(species.SpeciesBase)
	if err != nil {
		return newJSONError(err, http.StatusInternalServerError)
	}
	return nil
}

func genusIdFromName(genus_name string) (int64, error) {
	var genus_id struct{ Id int64 }
	q := `SELECT id FROM genera WHERE LOWER(genus_name) = LOWER($1);`
	if err := DBH.SelectOne(&genus_id, q, genus_name); err != nil {
		return 0, err
	}
	return genus_id.Id, nil
}

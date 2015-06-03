package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrSpeciesNotFound   = errors.New("species not found")
	ErrSpeciesNotUpdated = errors.New("species not updated")
)

func init() {
	DB.AddTableWithName(SpeciesBase{}, "species").SetKeys(true, "Id")
}

type SpeciesService struct{}

// SpeciesBase is what the DB expects to see for inserts/updates
type SpeciesBase struct {
	Id          int64      `db:"id" json:"id"`
	SpeciesName string     `db:"species_name" json:"speciesName"`
	TypeSpecies NullBool   `db:"type_species" json:"typeSpecies"`
	Etymology   NullString `db:"etymology" json:"etymology"`
	CreatedAt   NullTime   `db:"created_at" json:"createdAt"`
	UpdatedAt   NullTime   `db:"updated_at" json:"updatedAt"`
	DeletedAt   NullTime   `db:"deleted_at" json:"deletedAt"`
	CreatedBy   int64      `db:"created_by" json:"createdBy"`
	UpdatedBy   int64      `db:"updated_by" json:"updatedBy"`
	DeletedBy   NullInt64  `db:"deleted_by" json:"deletedBy"`
}

// Species & SpeciesJSON(s) are what ember expects to see
type Species struct {
	*SpeciesBase
	GenusName    string         `db:"genus_name" json:"genusName"`
	Strains      NullSliceInt64 `db:"strains" json:"strains"`
	TotalStrains int64          `db:"total_strains" json:"totalStrains"`
}

type ManySpecies []*Species

type SpeciesJSON struct {
	Species *Species `json:"species"`
}

type ManySpeciesJSON struct {
	ManySpecies *ManySpecies `json:"species"`
}

func (s *Species) marshal() ([]byte, error) {
	return json.Marshal(&SpeciesJSON{Species: s})
}

func (s *ManySpecies) marshal() ([]byte, error) {
	return json.Marshal(&ManySpeciesJSON{ManySpecies: s})
}

func (s SpeciesService) unmarshal(b []byte) (entity, error) {
	var sj SpeciesJSON
	err := json.Unmarshal(b, &sj)
	return sj.Species, err
}

func (s SpeciesService) list(opt *ListOptions) (entity, error) {
	if opt == nil {
		return nil, errors.New("must provide options")
	}
	var vals []interface{}

	sql := `SELECT sp.id, sp.species_name, sp.type_species, sp.etymology,
			sp.created_at, sp.created_by, sp.updated_at, sp.updated_by,
			sp.deleted_at, sp.deleted_by, g.genus_name, array_agg(st.id) AS strains,
			COUNT(st) AS total_strains
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

	var species ManySpecies
	err := DBH.Select(&species, sql, vals...)
	if err != nil {
		return nil, err
	}
	return &species, nil
}

func (s SpeciesService) get(id int64, genus string) (entity, error) {
	var species Species
	q := `SELECT sp.id, sp.species_name, sp.type_species, sp.etymology,
		sp.created_at, sp.created_by, sp.updated_at, sp.updated_by, sp.deleted_at,
		sp.deleted_by, g.genus_name, array_agg(st.id) AS strains,
		COUNT(st) AS total_strains
		FROM species sp
		INNER JOIN genera g ON g.id=sp.genus_id AND LOWER(g.genus_name)=$1
		LEFT OUTER JOIN strains st ON st.species_id=sp.id
		WHERE sp.id=$2
		GROUP BY sp.id, g.genus_name;`
	if err := DBH.SelectOne(&species, q, genus, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrStrainNotFound
		}
		return nil, err
	}
	return &species, nil
}

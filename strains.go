package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrStrainNotFound   = errors.New("strain not found")
	ErrStrainNotUpdated = errors.New("strain not updated")
)

func init() {
	DB.AddTableWithName(StrainBase{}, "strains").SetKeys(true, "Id")
}

type StrainService struct{}

// StrainBase is what the DB expects to see for inserts/updates
type StrainBase struct {
	Id               int64      `db:"id" json:"id"`
	SpeciesId        int64      `db:"species_id" json:"-"`
	StrainName       string     `db:"strain_name" json:"strainName"`
	TypeStrain       bool       `db:"type_strain" json:"typeStrain"`
	AccessionNumbers string     `db:"accession_numbers" json:"accessionNumbers"`
	Genbank          NullString `db:"genbank" json:"genbank"`
	IsolatedFrom     NullString `db:"isolated_from" json:"isolatedFrom"`
	Notes            NullString `db:"notes" json:"notes"`
	CreatedAt        NullTime   `db:"created_at" json:"createdAt"`
	UpdatedAt        NullTime   `db:"updated_at" json:"updatedAt"`
	DeletedAt        NullTime   `db:"deleted_at" json:"deletedAt"`
	CreatedBy        int64      `db:"created_by" json:"createdBy"`
	UpdatedBy        int64      `db:"updated_by" json:"updatedBy"`
	DeletedBy        NullInt64  `db:"deleted_by" json:"deletedBy"`
}

// Strain & StrainJSON(s) are what ember expects to see
type Strain struct {
	*StrainBase
	SpeciesName       string         `db:"species_name" json:"speciesName"`
	Measurements      NullSliceInt64 `db:"measurements" json:"measurements"`
	TotalMeasurements int64          `db:"total_measurements" json:"totalMeasurements"`
}

type Strains []*Strain

type StrainJSON struct {
	Strain *Strain `json:"strain"`
}

type StrainsJSON struct {
	Strains *Strains `json:"strains"`
}

func (s *Strain) marshal() ([]byte, error) {
	return json.Marshal(&StrainJSON{Strain: s})
}

func (s *Strains) marshal() ([]byte, error) {
	return json.Marshal(&StrainsJSON{Strains: s})
}

func (s StrainService) unmarshal(b []byte) (entity, error) {
	var sj StrainJSON
	err := json.Unmarshal(b, &sj)
	return sj.Strain, err
}

func (s StrainService) list(opt *ListOptions) (entity, error) {
	if opt == nil {
		return nil, errors.New("must provide options")
	}
	var vals []interface{}

	sql := `SELECT st.*, sp.species_name, array_agg(m.id) AS measurements,
		COUNT(m) AS total_measurements
		FROM strains st
		INNER JOIN species sp ON sp.id=st.species_id
		INNER JOIN genera g ON g.id=sp.genus_id AND LOWER(g.genus_name)=$1
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
		sql += " WHERE (" + strings.Join(conds, ") AND (") + ")"
	}

	sql += " GROUP BY st.id, sp.species_name;"

	var strains Strains
	err := DBH.Select(&strains, sql, vals...)
	if err != nil {
		return nil, err
	}
	return &strains, nil
}

func (s StrainService) get(id int64, genus string) (entity, error) {
	var strain Strain
	q := `SELECT st.*, sp.species_name, array_agg(m.id) AS measurements,
		COUNT(m) AS total_measurements
		FROM strains st
		INNER JOIN species sp ON sp.id=st.species_id
		INNER JOIN genera g ON g.id=sp.genus_id AND LOWER(g.genus_name)=$1
		LEFT OUTER JOIN measurements m ON m.strain_id=st.id
		WHERE st.id=$2
		GROUP BY st.id, sp.species_name;`
	if err := DBH.SelectOne(&strain, q, genus, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrStrainNotFound
		}
		return nil, err
	}
	return &strain, nil
}

func (s StrainService) update(id int64, e *entity, claims Claims) error {
	strain := (*e).(*Strain)
	strain.UpdatedBy = claims.Sub
	strain.UpdatedAt = currentTime()
	strain.Id = id

	var species_id struct{ Id int64 }
	q := `SELECT id FROM species WHERE species_name = $1;`
	if err := DBH.SelectOne(&species_id, q, strain.SpeciesName); err != nil {
		return err
	}
	strain.StrainBase.SpeciesId = species_id.Id

	count, err := DBH.Update(strain.StrainBase)
	if err != nil {
		return err
	}
	if count != 1 {
		return ErrStrainNotUpdated
	}
	return nil
}

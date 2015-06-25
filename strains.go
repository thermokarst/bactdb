package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	ErrStrainNotFound     = errors.New("Strain not found")
	ErrStrainNotFoundJSON = newJSONError(ErrStrainNotFound, http.StatusNotFound)
	ErrStrainNotUpdated   = errors.New("Strain not updated")
)

func init() {
	DB.AddTableWithName(StrainBase{}, "strains").SetKeys(true, "Id")
}

type StrainService struct{}

// StrainBase is what the DB expects to see for inserts/updates
type StrainBase struct {
	Id                  int64      `db:"id" json:"id"`
	SpeciesId           int64      `db:"species_id" json:"species,string"` // quirk in ember select
	StrainName          string     `db:"strain_name" json:"strainName"`
	TypeStrain          bool       `db:"type_strain" json:"typeStrain"`
	AccessionNumbers    string     `db:"accession_numbers" json:"accessionNumbers"`
	Genbank             NullString `db:"genbank" json:"genbank"`
	WholeGenomeSequence NullString `db:"whole_genome_sequence" json:"wholeGenomeSequence"`
	IsolatedFrom        NullString `db:"isolated_from" json:"isolatedFrom"`
	Notes               NullString `db:"notes" json:"notes"`
	CreatedAt           time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt           time.Time  `db:"updated_at" json:"updatedAt"`
	DeletedAt           NullTime   `db:"deleted_at" json:"deletedAt"`
	CreatedBy           int64      `db:"created_by" json:"createdBy"`
	UpdatedBy           int64      `db:"updated_by" json:"updatedBy"`
	DeletedBy           NullInt64  `db:"deleted_by" json:"deletedBy"`
}

// Strain & StrainJSON(s) are what ember expects to see
type Strain struct {
	*StrainBase
	Measurements      NullSliceInt64 `db:"measurements" json:"measurements"`
	TotalMeasurements int64          `db:"total_measurements" json:"totalMeasurements"`
	SortOrder         int64          `db:"sort_order" json:"sortOrder"`
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

func (s StrainService) list(val *url.Values) (entity, *appError) {
	if val == nil {
		return nil, ErrMustProvideOptionsJSON
	}
	var opt ListOptions
	if err := schemaDecoder.Decode(&opt, *val); err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	var vals []interface{}
	sql := `SELECT st.*, array_agg(m.id) AS measurements, COUNT(m) AS total_measurements,
		rank() OVER (ORDER BY sp.species_name ASC, st.type_strain ASC, st.strain_name ASC) AS sort_order
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

	sql += " GROUP BY st.id, st.species_id, sp.species_name;"

	strains := make(Strains, 0)
	err := DBH.Select(&strains, sql, vals...)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}
	return &strains, nil
}

func (s StrainService) get(id int64, genus string) (entity, *appError) {
	var strain Strain
	q := `SELECT st.*, array_agg(m.id) AS measurements, COUNT(m) AS total_measurements,
		0 AS sort_order
		FROM strains st
		INNER JOIN species sp ON sp.id=st.species_id
		INNER JOIN genera g ON g.id=sp.genus_id AND LOWER(g.genus_name)=$1
		LEFT OUTER JOIN measurements m ON m.strain_id=st.id
		WHERE st.id=$2
		GROUP BY st.id, st.species_id;`
	if err := DBH.SelectOne(&strain, q, genus, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrStrainNotFoundJSON
		}
		return nil, newJSONError(err, http.StatusInternalServerError)
	}
	return &strain, nil
}

func (s StrainService) update(id int64, e *entity, claims Claims) error {
	strain := (*e).(*Strain)
	strain.UpdatedBy = claims.Sub
	strain.UpdatedAt = time.Now()
	strain.Id = id

	count, err := DBH.Update(strain.StrainBase)
	if err != nil {
		return err
	}
	if count != 1 {
		return ErrStrainNotUpdated
	}
	return nil
}

func (s StrainService) create(e *entity, claims Claims) *appError {
	strain := (*e).(*Strain)
	ct := time.Now()
	strain.CreatedBy = claims.Sub
	strain.CreatedAt = ct
	strain.UpdatedBy = claims.Sub
	strain.UpdatedAt = ct

	if err := DBH.Insert(strain.StrainBase); err != nil {
		return newJSONError(err, http.StatusInternalServerError)
	}
	return nil
}

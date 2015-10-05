package models

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/thermokarst/bactdb/Godeps/_workspace/src/github.com/jmoiron/modl"
	"github.com/thermokarst/bactdb/errors"
	"github.com/thermokarst/bactdb/helpers"
	"github.com/thermokarst/bactdb/types"
)

func init() {
	DB.AddTableWithName(StrainBase{}, "strains").SetKeys(true, "ID")
}

// PreInsert is a modl hook.
func (s *StrainBase) PreInsert(e modl.SqlExecutor) error {
	ct := helpers.CurrentTime()
	s.CreatedAt = ct
	s.UpdatedAt = ct
	return nil
}

// PreUpdate is a modl hook.
func (s *StrainBase) PreUpdate(e modl.SqlExecutor) error {
	s.UpdatedAt = helpers.CurrentTime()
	return nil
}

// StrainBase is what the DB expects for write operations.
type StrainBase struct {
	ID                  int64            `db:"id" json:"id"`
	SpeciesID           int64            `db:"species_id" json:"species"`
	StrainName          string           `db:"strain_name" json:"strainName"`
	TypeStrain          bool             `db:"type_strain" json:"typeStrain"`
	AccessionNumbers    types.NullString `db:"accession_numbers" json:"accessionNumbers"`
	Genbank             types.NullString `db:"genbank" json:"genbank"`
	WholeGenomeSequence types.NullString `db:"whole_genome_sequence" json:"wholeGenomeSequence"`
	IsolatedFrom        types.NullString `db:"isolated_from" json:"isolatedFrom"`
	Notes               types.NullString `db:"notes" json:"notes"`
	CreatedAt           types.NullTime   `db:"created_at" json:"createdAt"`
	UpdatedAt           types.NullTime   `db:"updated_at" json:"updatedAt"`
	DeletedAt           types.NullTime   `db:"deleted_at" json:"deletedAt"`
	CreatedBy           int64            `db:"created_by" json:"createdBy"`
	UpdatedBy           int64            `db:"updated_by" json:"updatedBy"`
	DeletedBy           types.NullInt64  `db:"deleted_by" json:"deletedBy"`
}

// Strain is what the DB expects for read operations, and is what the API expects
// to return to the requester.
type Strain struct {
	*StrainBase
	Measurements      types.NullSliceInt64 `db:"measurements" json:"measurements"`
	Characteristics   types.NullSliceInt64 `db:"characteristics" json:"characteristics"`
	TotalMeasurements int64                `db:"total_measurements" json:"totalMeasurements"`
	SortOrder         int64                `db:"sort_order" json:"sortOrder"`
	CanEdit           bool                 `db:"-" json:"canEdit"`
}

// Strains are multiple strain entities.
type Strains []*Strain

// StrainMeta stashes some metadata related to the entity.
type StrainMeta struct {
	CanAdd bool `json:"canAdd"`
}

// SpeciesName returns a strain's species name.
func (s StrainBase) SpeciesName() string {
	var species SpeciesBase
	if err := DBH.Get(&species, s.SpeciesID); err != nil {
		return ""
	}
	return species.SpeciesName
}

// ListStrains returns all strains.
func ListStrains(opt helpers.ListOptions, claims *types.Claims) (*Strains, error) {
	var vals []interface{}

	q := `SELECT st.*, array_agg(m.id) AS measurements,
		array_agg(DISTINCT m.characteristic_id) AS characteristics,
		COUNT(m) AS total_measurements,
		rank() OVER (ORDER BY sp.species_name ASC, st.type_strain ASC, st.strain_name ASC) AS sort_order
		FROM strains st
		INNER JOIN species sp ON sp.id=st.species_id
		INNER JOIN genera g ON g.id=sp.genus_id AND LOWER(g.genus_name)=LOWER($1)
		LEFT OUTER JOIN measurements m ON m.strain_id=st.id`
	vals = append(vals, opt.Genus)

	if len(opt.IDs) != 0 {
		var conds []string
		s := "st.id IN ("
		for i, id := range opt.IDs {
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

	for _, s := range strains {
		s.CanEdit = helpers.CanEdit(claims, s.CreatedBy)
	}

	return &strains, nil
}

// GetStrain returns a particular strain.
func GetStrain(id int64, genus string, claims *types.Claims) (*Strain, error) {
	var strain Strain
	q := `SELECT st.*, array_agg(DISTINCT m.id) AS measurements,
		array_agg(DISTINCT m.characteristic_id) AS characteristics,
		COUNT(m) AS total_measurements, 0 AS sort_order
		FROM strains st
		INNER JOIN species sp ON sp.id=st.species_id
		INNER JOIN genera g ON g.id=sp.genus_id AND LOWER(g.genus_name)=LOWER($1)
		LEFT OUTER JOIN measurements m ON m.strain_id=st.id
		WHERE st.id=$2
		GROUP BY st.id;`
	if err := DBH.SelectOne(&strain, q, genus, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrStrainNotFound
		}
		return nil, err
	}

	strain.CanEdit = helpers.CanEdit(claims, strain.CreatedBy)

	return &strain, nil
}

// SpeciesOptsFromStrains returns the options for finding all related species for a
// set of strains.
func SpeciesOptsFromStrains(opt helpers.ListOptions) (*helpers.ListOptions, error) {
	var relatedSpeciesIDs []int64

	if opt.IDs == nil || len(opt.IDs) == 0 {
		q := `SELECT DISTINCT st.species_id
			FROM strains st
			INNER JOIN species sp ON sp.id=st.species_id
			INNER JOIN genera g ON g.id=sp.genus_id AND LOWER(g.genus_name)=LOWER($1);`
		if err := DBH.Select(&relatedSpeciesIDs, q, opt.Genus); err != nil {
			return nil, err
		}
	} else {
		var vals []interface{}
		var count int64 = 1
		q := fmt.Sprintf("SELECT DISTINCT species_id FROM strains WHERE %s;", helpers.ValsIn("id", opt.IDs, &vals, &count))
		if err := DBH.Select(&relatedSpeciesIDs, q, vals...); err != nil {
			return nil, err
		}
	}

	return &helpers.ListOptions{Genus: opt.Genus, IDs: relatedSpeciesIDs}, nil
}

// CharacteristicsOptsFromStrains returns the options for finding all related
// characteristics for a set of strains.
func CharacteristicsOptsFromStrains(opt helpers.ListOptions) (*helpers.ListOptions, error) {
	var relatedCharacteristicsIDs []int64

	if opt.IDs == nil || len(opt.IDs) == 0 {
		q := `SELECT DISTINCT m.characteristic_id
				FROM measurements m
				INNER JOIN strains st ON st.id=m.strain_id
				INNER JOIN species sp ON sp.id=st.species_id
				INNER JOIN genera g ON g.id=sp.genus_id AND LOWER(g.genus_name)=LOWER($1);`
		if err := DBH.Select(&relatedCharacteristicsIDs, q, opt.Genus); err != nil {
			return nil, err
		}
	} else {
		var vals []interface{}
		var count int64 = 1
		q := fmt.Sprintf("SELECT DISTINCT characteristic_id FROM measurements WHERE %s;", helpers.ValsIn("strain_id", opt.IDs, &vals, &count))
		if err := DBH.Select(&relatedCharacteristicsIDs, q, vals...); err != nil {
			return nil, err
		}
	}

	return &helpers.ListOptions{Genus: opt.Genus, IDs: relatedCharacteristicsIDs}, nil
}

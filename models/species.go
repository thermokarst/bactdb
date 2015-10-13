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
	DB.AddTableWithName(SpeciesBase{}, "species").SetKeys(true, "ID")
}

// PreInsert is a modl hook.
func (s *SpeciesBase) PreInsert(e modl.SqlExecutor) error {
	ct := helpers.CurrentTime()
	s.CreatedAt = ct
	s.UpdatedAt = ct
	return nil
}

// PreUpdate is a modl hook.
func (s *SpeciesBase) PreUpdate(e modl.SqlExecutor) error {
	s.UpdatedAt = helpers.CurrentTime()
	return nil
}

// UpdateError satisfies base interface.
func (s *SpeciesBase) UpdateError() error {
	return errors.ErrSpeciesNotUpdated
}

// DeleteError satisfies base interface.
func (s *SpeciesBase) DeleteError() error {
	return errors.ErrSpeciesNotDeleted
}

// SpeciesBase is what the DB expects for write operations.
type SpeciesBase struct {
	ID                  int64            `db:"id" json:"id"`
	GenusID             int64            `db:"genus_id" json:"-"`
	SubspeciesSpeciesID types.NullInt64  `db:"subspecies_species_id" json:"-"`
	SpeciesName         string           `db:"species_name" json:"speciesName"`
	TypeSpecies         types.NullBool   `db:"type_species" json:"typeSpecies"`
	Etymology           types.NullString `db:"etymology" json:"etymology"`
	CreatedAt           types.NullTime   `db:"created_at" json:"createdAt"`
	UpdatedAt           types.NullTime   `db:"updated_at" json:"updatedAt"`
	DeletedAt           types.NullTime   `db:"deleted_at" json:"deletedAt"`
	CreatedBy           int64            `db:"created_by" json:"createdBy"`
	UpdatedBy           int64            `db:"updated_by" json:"updatedBy"`
	DeletedBy           types.NullInt64  `db:"deleted_by" json:"deletedBy"`
}

// Species is what the DB expects for read operations, and is what the API expects
// to return to the requester.
type Species struct {
	*SpeciesBase
	GenusName    string               `db:"genus_name" json:"genusName"`
	Strains      types.NullSliceInt64 `db:"strains" json:"strains"`
	TotalStrains int64                `db:"total_strains" json:"totalStrains"`
	SortOrder    int64                `db:"sort_order" json:"sortOrder"`
	CanEdit      bool                 `db:"-" json:"canEdit"`
}

// ManySpecies is multiple species entities.
type ManySpecies []*Species

// SpeciesMeta stashes some metadata related to the entity.
type SpeciesMeta struct {
	CanAdd bool `json:"canAdd"`
}

// GenusIDFromName looks up the genus' ID.
func GenusIDFromName(genusName string) (int64, error) {
	var genusID struct{ ID int64 }
	q := `SELECT id FROM genera WHERE LOWER(genus_name) = LOWER($1);`
	if err := DBH.SelectOne(&genusID, q, genusName); err != nil {
		return 0, err
	}
	return genusID.ID, nil
}

// StrainOptsFromSpecies returns the options for finding all related strains for
// a set of species.
func StrainOptsFromSpecies(opt helpers.ListOptions) (*helpers.ListOptions, error) {
	var relatedStrainIDs []int64

	if opt.IDs == nil {
		q := `SELECT DISTINCT st.id
			FROM strains st
			INNER JOIN species sp ON sp.id=st.species_id
			INNER JOIN genera g ON g.id=sp.genus_id AND LOWER(g.genus_name)=LOWER($1);`
		if err := DBH.Select(&relatedStrainIDs, q, opt.Genus); err != nil {
			return nil, err
		}
	} else {
		var vals []interface{}
		var count int64 = 1
		q := fmt.Sprintf("SELECT DISTINCT id FROM strains WHERE %s;", helpers.ValsIn("species_id", opt.IDs, &vals, &count))

		if err := DBH.Select(&relatedStrainIDs, q, vals...); err != nil {
			return nil, err
		}
	}

	return &helpers.ListOptions{Genus: opt.Genus, IDs: relatedStrainIDs}, nil
}

// StrainsFromSpeciesID returns the options for finding all related strains for a
// particular species.
func StrainsFromSpeciesID(id int64, genus string, claims *types.Claims) (*Strains, error) {
	opt := helpers.ListOptions{
		Genus: genus,
		IDs:   []int64{id},
	}

	strainsOpt, err := StrainOptsFromSpecies(opt)
	if err != nil {
		return nil, err
	}

	strains, err := ListStrains(*strainsOpt, claims)
	if err != nil {
		return nil, err
	}

	return strains, nil
}

// ListSpecies returns all species
func ListSpecies(opt helpers.ListOptions, claims *types.Claims) (*ManySpecies, error) {
	var vals []interface{}

	q := `SELECT sp.*, g.genus_name, array_agg(st.id) AS strains,
			COUNT(st) AS total_strains,
			rank() OVER (ORDER BY sp.species_name ASC) AS sort_order
			FROM species sp
			INNER JOIN genera g ON g.id=sp.genus_id AND LOWER(g.genus_name)=LOWER($1)
			LEFT OUTER JOIN strains st ON st.species_id=sp.id`
	vals = append(vals, opt.Genus)

	if len(opt.IDs) != 0 {
		var conds []string
		s := "sp.id IN ("
		for i, id := range opt.IDs {
			s = s + fmt.Sprintf("$%v,", i+2) // start param index at 2
			vals = append(vals, id)
		}
		s = s[:len(s)-1] + ")"
		conds = append(conds, s)
		q += " WHERE (" + strings.Join(conds, ") AND (") + ")"
	}

	q += " GROUP BY sp.id, g.genus_name;"

	species := make(ManySpecies, 0)
	err := DBH.Select(&species, q, vals...)
	if err != nil {
		return nil, err
	}

	for _, s := range species {
		s.CanEdit = helpers.CanEdit(claims, s.CreatedBy)
	}

	return &species, nil
}

// GetSpecies returns a particular species.
func GetSpecies(id int64, genus string, claims *types.Claims) (*Species, error) {
	var species Species
	q := `SELECT sp.*, g.genus_name, array_agg(st.id) AS strains,
		COUNT(st) AS total_strains, 0 AS sort_order
		FROM species sp
		INNER JOIN genera g ON g.id=sp.genus_id AND LOWER(g.genus_name)=LOWER($1)
		LEFT OUTER JOIN strains st ON st.species_id=sp.id
		WHERE sp.id=$2
		GROUP BY sp.id, g.genus_name;`
	if err := DBH.SelectOne(&species, q, genus, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrSpeciesNotFound
		}
		return nil, err
	}

	species.CanEdit = helpers.CanEdit(claims, species.CreatedBy)

	return &species, nil
}

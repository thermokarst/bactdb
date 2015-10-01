package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/thermokarst/bactdb/Godeps/_workspace/src/github.com/jmoiron/modl"
	"github.com/thermokarst/bactdb/helpers"
	"github.com/thermokarst/bactdb/types"
)

var (
	ErrSpeciesNotFound   = errors.New("Species not found")
	ErrSpeciesNotUpdated = errors.New("Species not updated")
)

func init() {
	DB.AddTableWithName(SpeciesBase{}, "species").SetKeys(true, "Id")
}

func (s *SpeciesBase) PreInsert(e modl.SqlExecutor) error {
	ct := helpers.CurrentTime()
	s.CreatedAt = ct
	s.UpdatedAt = ct
	return nil
}

func (s *SpeciesBase) PreUpdate(e modl.SqlExecutor) error {
	s.UpdatedAt = helpers.CurrentTime()
	return nil
}

type SpeciesBase struct {
	Id                  int64            `db:"id" json:"id"`
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

type Species struct {
	*SpeciesBase
	GenusName    string               `db:"genus_name" json:"genusName"`
	Strains      types.NullSliceInt64 `db:"strains" json:"strains"`
	TotalStrains int64                `db:"total_strains" json:"totalStrains"`
	SortOrder    int64                `db:"sort_order" json:"sortOrder"`
	CanEdit      bool                 `db:"-" json:"canEdit"`
}

type ManySpecies []*Species

type SpeciesMeta struct {
	CanAdd bool `json:"canAdd"`
}

func GenusIdFromName(genus_name string) (int64, error) {
	var genus_id struct{ Id int64 }
	q := `SELECT id FROM genera WHERE LOWER(genus_name) = LOWER($1);`
	if err := DBH.SelectOne(&genus_id, q, genus_name); err != nil {
		return 0, err
	}
	return genus_id.Id, nil
}

func StrainOptsFromSpecies(opt helpers.ListOptions) (*helpers.ListOptions, error) {
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
		q := fmt.Sprintf("SELECT DISTINCT id FROM strains WHERE %s;", helpers.ValsIn("species_id", opt.Ids, &vals, &count))

		if err := DBH.Select(&relatedStrainIds, q, vals...); err != nil {
			return nil, err
		}
	}

	return &helpers.ListOptions{Genus: opt.Genus, Ids: relatedStrainIds}, nil
}

func StrainsFromSpeciesId(id int64, genus string, claims *types.Claims) (*Strains, error) {
	opt := helpers.ListOptions{
		Genus: genus,
		Ids:   []int64{id},
	}

	strains_opt, err := StrainOptsFromSpecies(opt)
	if err != nil {
		return nil, err
	}

	strains, err := ListStrains(*strains_opt, claims)
	if err != nil {
		return nil, err
	}

	return strains, nil
}

func ListSpecies(opt helpers.ListOptions, claims *types.Claims) (*ManySpecies, error) {
	var vals []interface{}

	q := `SELECT sp.*, g.genus_name, array_agg(st.id) AS strains,
			COUNT(st) AS total_strains,
			rank() OVER (ORDER BY sp.species_name ASC) AS sort_order
			FROM species sp
			INNER JOIN genera g ON g.id=sp.genus_id AND LOWER(g.genus_name)=LOWER($1)
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
			return nil, ErrSpeciesNotFound
		}
		return nil, err
	}

	species.CanEdit = helpers.CanEdit(claims, species.CreatedBy)

	return &species, nil
}

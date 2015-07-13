package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/jmoiron/modl"
)

var (
	ErrCharacteristicNotFound     = errors.New("Characteristic not found")
	ErrCharacteristicNotFoundJSON = newJSONError(ErrCharacteristicNotFound, http.StatusNotFound)
)

func init() {
	DB.AddTableWithName(CharacteristicBase{}, "characteristics").SetKeys(true, "Id")
}

func (c *CharacteristicBase) PreInsert(e modl.SqlExecutor) error {
	ct := currentTime()
	c.CreatedAt = ct
	c.UpdatedAt = ct
	return nil
}

func (c *CharacteristicBase) PreUpdate(e modl.SqlExecutor) error {
	c.UpdatedAt = currentTime()
	return nil
}

type CharacteristicService struct{}

type CharacteristicBase struct {
	Id                   int64     `json:"id,omitempty"`
	CharacteristicName   string    `db:"characteristic_name" json:"characteristicName"`
	CharacteristicTypeId int64     `db:"characteristic_type_id" json:"-"`
	SortOrder            NullInt64 `db:"sort_order" json:"sortOrder"`
	CreatedAt            NullTime  `db:"created_at" json:"createdAt"`
	UpdatedAt            NullTime  `db:"updated_at" json:"updatedAt"`
	DeletedAt            NullTime  `db:"deleted_at" json:"deletedAt"`
	CreatedBy            int64     `db:"created_by" json:"createdBy"`
	UpdatedBy            int64     `db:"updated_by" json:"updatedBy"`
	DeletedBy            NullInt64 `db:"deleted_by" json:"deletedBy"`
}

type Characteristic struct {
	*CharacteristicBase
	Measurements       NullSliceInt64 `db:"measurements" json:"measurements"`
	Strains            NullSliceInt64 `db:"strains" json:"strains"`
	CharacteristicType string         `db:"characteristic_type_name" json:"characteristicTypeName"`
	CanEdit            bool           `db:"-" json:"canEdit"`
}

type Characteristics []*Characteristic

type CharacteristicMeta struct {
	CanAdd bool `json:"canAdd"`
}

type CharacteristicPayload struct {
	Characteristic *Characteristic     `json:"characteristic"`
	Measurements   *Measurements       `json:"measurements"`
	Strains        *Strains            `json:"strains"`
	Species        *ManySpecies        `json:"species"`
	Meta           *CharacteristicMeta `json:"meta"`
}

type CharacteristicsPayload struct {
	Characteristics *Characteristics    `json:"characteristics"`
	Measurements    *Measurements       `json:"measurements"`
	Strains         *Strains            `json:"strains"`
	Species         *ManySpecies        `json:"species"`
	Meta            *CharacteristicMeta `json:"meta"`
}

func (c *CharacteristicPayload) marshal() ([]byte, error) {
	return json.Marshal(c)
}

func (c *CharacteristicsPayload) marshal() ([]byte, error) {
	return json.Marshal(c)
}

func (c CharacteristicService) list(val *url.Values, claims *Claims) (entity, *appError) {
	if val == nil {
		return nil, ErrMustProvideOptionsJSON
	}
	var opt ListOptions
	if err := schemaDecoder.Decode(&opt, *val); err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	characteristics, err := listCharacteristics(opt, claims)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	strains_opt, err := strainOptsFromCharacteristics(opt)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	strains, err := listStrains(*strains_opt, claims)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	species_opt, err := speciesOptsFromStrains(*strains_opt)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	species, err := listSpecies(*species_opt, claims)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	// TODO: tack on measurements
	payload := CharacteristicsPayload{
		Characteristics: characteristics,
		Measurements:    nil,
		Strains:         strains,
		Species:         species,
		Meta: &CharacteristicMeta{
			CanAdd: canAdd(claims),
		},
	}

	return &payload, nil
}

func (c CharacteristicService) get(id int64, genus string, claims *Claims) (entity, *appError) {
	characteristic, err := getCharacteristic(id, claims)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	strains, strain_opts, err := strainsFromCharacteristicId(id, genus, claims)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	species_opt, err := speciesOptsFromStrains(*strain_opts)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	species, err := listSpecies(*species_opt, claims)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	// TODO: tack on measurements
	payload := CharacteristicPayload{
		Characteristic: characteristic,
		Measurements:   nil,
		Strains:        strains,
		Species:        species,
		Meta: &CharacteristicMeta{
			CanAdd: canAdd(claims),
		},
	}

	return &payload, nil
}

func listCharacteristics(opt ListOptions, claims *Claims) (*Characteristics, error) {
	var vals []interface{}

	q := `SELECT c.*, array_agg(m.id) AS measurements,
			array_agg(st.id) AS strains, ct.characteristic_type_name
			FROM characteristics c
			INNER JOIN characteristic_types ct ON ct.id=c.characteristic_type_id
			LEFT OUTER JOIN measurements m ON m.characteristic_id=c.id
			LEFT OUTER JOIN strains st ON st.id=m.strain_id`

	if len(opt.Ids) != 0 {
		var counter int64 = 1
		w := valsIn("c.id", opt.Ids, &vals, &counter)

		q += fmt.Sprintf(" WHERE %s", w)
	}

	q += " GROUP BY c.id, ct.characteristic_type_name;"

	characteristics := make(Characteristics, 0)
	err := DBH.Select(&characteristics, q, vals...)
	if err != nil {
		return nil, err
	}

	for _, c := range characteristics {
		c.CanEdit = canEdit(claims, c.CreatedBy)
	}

	return &characteristics, nil
}

func strainOptsFromCharacteristics(opt ListOptions) (*ListOptions, error) {
	relatedStrainIds := make([]int64, 0)

	if opt.Ids == nil {
		q := `SELECT DISTINCT strain_id FROM measurements;`
		if err := DBH.Select(&relatedStrainIds, q); err != nil {
			return nil, err
		}
	} else {
		var vals []interface{}
		var count int64 = 1
		q := fmt.Sprintf("SELECT DISTINCT strain_id FROM measurements WHERE %s;", valsIn("characteristic_id", opt.Ids, &vals, &count))

		if err := DBH.Select(&relatedStrainIds, q, vals...); err != nil {
			return nil, err
		}
	}

	return &ListOptions{Genus: opt.Genus, Ids: relatedStrainIds}, nil
}

func strainsFromCharacteristicId(id int64, genus string, claims *Claims) (*Strains, *ListOptions, error) {
	opt := ListOptions{
		Genus: genus,
		Ids:   []int64{id},
	}

	strains_opt, err := strainOptsFromCharacteristics(opt)
	if err != nil {
		return nil, nil, err
	}

	strains, err := listStrains(*strains_opt, claims)
	if err != nil {
		return nil, nil, err
	}

	return strains, strains_opt, nil
}

func getCharacteristic(id int64, claims *Claims) (*Characteristic, error) {
	var characteristic Characteristic
	q := `SELECT c.*, array_agg(m.id) AS measurements,
			array_agg(st.id) AS strains, ct.characteristic_type_name
			FROM characteristics c
			INNER JOIN characteristic_types ct ON ct.id=c.characteristic_type_id
			LEFT OUTER JOIN measurements m ON m.characteristic_id=c.id
			LEFT OUTER JOIN strains st ON st.id=m.strain_id
			WHERE c.id=$1
			GROUP BY c.id, ct.characteristic_type_name;`
	if err := DBH.SelectOne(&characteristic, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrCharacteristicNotFound
		}
		return nil, err
	}

	characteristic.CanEdit = canEdit(claims, characteristic.CreatedBy)

	return &characteristic, nil
}

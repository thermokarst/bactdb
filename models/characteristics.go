package models

import (
	"database/sql"
	"fmt"

	"github.com/thermokarst/bactdb/Godeps/_workspace/src/github.com/jmoiron/modl"
	"github.com/thermokarst/bactdb/errors"
	"github.com/thermokarst/bactdb/helpers"
	"github.com/thermokarst/bactdb/types"
)

func init() {
	DB.AddTableWithName(CharacteristicBase{}, "characteristics").SetKeys(true, "ID")
}

// PreInsert is a modl hook
func (c *CharacteristicBase) PreInsert(e modl.SqlExecutor) error {
	ct := helpers.CurrentTime()
	c.CreatedAt = ct
	c.UpdatedAt = ct
	return nil
}

// PreUpdate is a modl hook
func (c *CharacteristicBase) PreUpdate(e modl.SqlExecutor) error {
	c.UpdatedAt = helpers.CurrentTime()
	return nil
}

// UpdateError satisfies base interface.
func (c *CharacteristicBase) UpdateError() error {
	return errors.ErrCharacteristicNotUpdated
}

// DeleteError satisfies base interface.
func (c *CharacteristicBase) DeleteError() error {
	return errors.ErrCharacteristicNotDeleted
}

func (c *CharacteristicBase) validate() types.ValidationError {
	cv := make(types.ValidationError, 0)

	if c.CharacteristicName == "" {
		cv = append(cv, types.NewValidationError(
			"characteristicName",
			helpers.MustProvideAValue))
	}

	if c.CharacteristicTypeID == 0 {
		cv = append(cv, types.NewValidationError(
			"characteristicType",
			helpers.MustProvideAValue))
	}

	if len(cv) > 0 {
		return cv
	}

	return nil
}

// CharacteristicBase is what the DB expects for write operations
type CharacteristicBase struct {
	ID                   int64           `json:"id,omitempty"`
	CharacteristicName   string          `db:"characteristic_name" json:"characteristicName"`
	CharacteristicTypeID int64           `db:"characteristic_type_id" json:"-"`
	SortOrder            types.NullInt64 `db:"sort_order" json:"sortOrder"`
	CreatedAt            types.NullTime  `db:"created_at" json:"createdAt"`
	UpdatedAt            types.NullTime  `db:"updated_at" json:"updatedAt"`
	CreatedBy            int64           `db:"created_by" json:"createdBy"`
	UpdatedBy            int64           `db:"updated_by" json:"updatedBy"`
}

// Characteristic is what the DB expects for read operations, and is what the API
// expects to return to the requester.
type Characteristic struct {
	*CharacteristicBase
	Measurements       types.NullSliceInt64 `db:"measurements" json:"measurements"`
	Strains            types.NullSliceInt64 `db:"strains" json:"strains"`
	CharacteristicType string               `db:"characteristic_type_name" json:"characteristicTypeName"`
	CanEdit            bool                 `db:"-" json:"canEdit"`
}

// Characteristics are multiple characteristic entities
type Characteristics []*Characteristic

// ListCharacteristics returns all characteristics
func ListCharacteristics(opt helpers.ListOptions, claims *types.Claims) (*Characteristics, error) {
	var vals []interface{}

	q := `SELECT c.*, ct.characteristic_type_name,
			array_agg(DISTINCT st.id) AS strains, array_agg(DISTINCT m.id) AS measurements
			FROM strains st
			INNER JOIN species sp ON sp.id=st.species_id
			INNER JOIN genera g ON g.id=sp.genus_id AND LOWER(g.genus_name)=LOWER($1)
			INNER JOIN measurements m ON m.strain_id=st.id
			RIGHT OUTER JOIN characteristics c ON c.id=m.characteristic_id
			INNER JOIN characteristic_types ct ON ct.id=c.characteristic_type_id`
	vals = append(vals, opt.Genus)

	if len(opt.IDs) != 0 {
		var counter int64 = 2
		w := helpers.ValsIn("c.id", opt.IDs, &vals, &counter)

		q += fmt.Sprintf(" WHERE %s", w)
	}

	q += ` GROUP BY c.id, ct.characteristic_type_name
			ORDER BY ct.characteristic_type_name, c.sort_order ASC;`

	var characteristics Characteristics
	err := DBH.Select(&characteristics, q, vals...)
	if err != nil {
		return nil, err
	}

	for _, c := range characteristics {
		c.CanEdit = helpers.CanEdit(claims, c.CreatedBy)
	}

	return &characteristics, nil
}

// StrainOptsFromCharacteristics returns the options for finding all related strains
// for a set of characteristics.
func StrainOptsFromCharacteristics(opt helpers.ListOptions) (*helpers.ListOptions, error) {
	var relatedStrainIDs []int64
	baseQ := `SELECT DISTINCT m.strain_id
		FROM measurements m
		INNER JOIN strains st ON st.id=m.strain_id
		INNER JOIN species sp ON sp.id=st.species_id
		INNER JOIN genera g ON g.id=sp.genus_id AND LOWER(g.genus_name)=LOWER($1)`
	if opt.IDs == nil {
		q := fmt.Sprintf("%s;", baseQ)
		if err := DBH.Select(&relatedStrainIDs, q, opt.Genus); err != nil {
			return nil, err
		}
	} else {
		var vals []interface{}
		var count int64 = 2
		vals = append(vals, opt.Genus)
		q := fmt.Sprintf("%s WHERE %s ", baseQ, helpers.ValsIn("m.characteristic_id", opt.IDs, &vals, &count))

		if err := DBH.Select(&relatedStrainIDs, q, vals...); err != nil {
			return nil, err
		}
	}

	return &helpers.ListOptions{Genus: opt.Genus, IDs: relatedStrainIDs}, nil
}

// MeasurementOptsFromCharacteristics returns the options for finding all related
// measurements for a set of characteristics.
func MeasurementOptsFromCharacteristics(opt helpers.ListOptions) (*helpers.MeasurementListOptions, error) {
	var relatedMeasurementIDs []int64
	baseQ := `SELECT m.id
		FROM measurements m
		INNER JOIN strains st ON st.id=m.strain_id
		INNER JOIN species sp ON sp.id=st.species_id
		INNER JOIN genera g ON g.id=sp.genus_id AND LOWER(g.genus_name)=LOWER($1)`

	if opt.IDs == nil {
		q := fmt.Sprintf("%s;", baseQ)
		if err := DBH.Select(&relatedMeasurementIDs, q, opt.Genus); err != nil {
			return nil, err
		}
	} else {
		var vals []interface{}
		var count int64 = 2
		vals = append(vals, opt.Genus)
		q := fmt.Sprintf("%s WHERE %s;", baseQ, helpers.ValsIn("characteristic_id", opt.IDs, &vals, &count))

		if err := DBH.Select(&relatedMeasurementIDs, q, vals...); err != nil {
			return nil, err
		}
	}

	return &helpers.MeasurementListOptions{ListOptions: helpers.ListOptions{Genus: opt.Genus, IDs: relatedMeasurementIDs}, Strains: nil, Characteristics: nil}, nil
}

// StrainsFromCharacteristicID returns a set of strains (as well as the options for
// finding those strains) for a particular characteristic.
func StrainsFromCharacteristicID(id int64, genus string, claims *types.Claims) (*Strains, *helpers.ListOptions, error) {
	opt := helpers.ListOptions{
		Genus: genus,
		IDs:   []int64{id},
	}

	strainsOpt, err := StrainOptsFromCharacteristics(opt)
	if err != nil {
		return nil, nil, err
	}

	strains, err := ListStrains(*strainsOpt, claims)
	if err != nil {
		return nil, nil, err
	}

	return strains, strainsOpt, nil
}

// MeasurementsFromCharacteristicID returns a set of measurements (as well as the
// options for finding those measurements) for a particular characteristic.
func MeasurementsFromCharacteristicID(id int64, genus string, claims *types.Claims) (*Measurements, *helpers.MeasurementListOptions, error) {
	opt := helpers.ListOptions{
		Genus: genus,
		IDs:   []int64{id},
	}

	measurementOpt, err := MeasurementOptsFromCharacteristics(opt)
	if err != nil {
		return nil, nil, err
	}

	measurements, err := ListMeasurements(*measurementOpt, claims)
	if err != nil {
		return nil, nil, err
	}

	return measurements, measurementOpt, nil
}

// GetCharacteristic returns a particular characteristic.
func GetCharacteristic(id int64, genus string, claims *types.Claims) (*Characteristic, error) {
	var characteristic Characteristic
	q := `SELECT c.*, ct.characteristic_type_name,
		array_agg(DISTINCT st.id) AS strains, array_agg(DISTINCT m.id) AS measurements
		FROM strains st
		INNER JOIN species sp ON sp.id=st.species_id
		INNER JOIN genera g ON g.id=sp.genus_id AND LOWER(g.genus_name)=LOWER($1)
		INNER JOIN measurements m ON m.strain_id=st.id
		RIGHT OUTER JOIN characteristics c ON c.id=m.characteristic_id
		INNER JOIN characteristic_types ct ON ct.id=c.characteristic_type_id
		WHERE c.id=$2
		GROUP BY c.id, ct.characteristic_type_name;`
	if err := DBH.SelectOne(&characteristic, q, genus, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrCharacteristicNotFound
		}
		return nil, err
	}

	characteristic.CanEdit = helpers.CanEdit(claims, characteristic.CreatedBy)

	return &characteristic, nil
}

// InsertOrGetCharacteristicType performs an UPSERT operation on the database
// for a characteristic type.
func InsertOrGetCharacteristicType(val string, claims *types.Claims) (int64, error) {
	var id int64
	q := `SELECT id FROM characteristic_types WHERE characteristic_type_name=$1;`
	if err := DBH.SelectOne(&id, q, val); err != nil {
		if err == sql.ErrNoRows {
			i := `INSERT INTO characteristic_types
				(characteristic_type_name, created_at, updated_at, created_by, updated_by)
				VALUES ($1, $2, $3, $4, $5) RETURNING id;`
			ct := helpers.CurrentTime()
			if err := DB.Db.QueryRow(i, val, ct, ct, claims.Sub, claims.Sub).Scan(&id); err != nil {
				return 0, err
			}
		} else {
			return 0, err
		}
	}
	return id, nil
}

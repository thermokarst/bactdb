package models

import (
	"github.com/thermokarst/bactdb/Godeps/_workspace/src/github.com/jmoiron/modl"
	"github.com/thermokarst/bactdb/types"
)

type base interface {
	PreInsert(modl.SqlExecutor) error
	PreUpdate(modl.SqlExecutor) error
	UpdateError() error
	DeleteError() error
	validate() types.ValidationError
}

// Create will create a new DB record of a model.
func Create(b base) error {
	if err := b.validate(); err != nil {
		return err
	}

	if err := DBH.Insert(b); err != nil {
		return nil
	}
	return nil
}

// Update runs a DB update on a model.
func Update(b base) error {
	if err := b.validate(); err != nil {
		return err
	}

	count, err := DBH.Update(b)
	if err != nil {
		return err
	}
	if count != 1 {
		return b.UpdateError()
	}

	return nil
}

// Delete runs a DB delete on a model.
func Delete(b base) error {
	count, err := DBH.Delete(b)
	if err != nil {
		return err
	}
	if count != 1 {
		return b.DeleteError()
	}

	return nil
}

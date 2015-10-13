package models

import "github.com/thermokarst/bactdb/Godeps/_workspace/src/github.com/jmoiron/modl"

type base interface {
	PreInsert(modl.SqlExecutor) error
	PreUpdate(modl.SqlExecutor) error
	UpdateError() error
}

// Create will create a new DB record of a model.
func Create(b base) error {
	if err := DBH.Insert(b); err != nil {
		return nil
	}
	return nil
}

// Update runs a DB update on a model.
func Update(b base) error {
	count, err := DBH.Update(b)
	if err != nil {
		return err
	}
	if count != 1 {
		return b.UpdateError()
	}
	return nil
}

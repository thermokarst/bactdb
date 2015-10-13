package models

import "github.com/thermokarst/bactdb/Godeps/_workspace/src/github.com/jmoiron/modl"

type base interface {
	PreCreate(modl.SqlExecutor) error
	PreUpdate(modl.SqlExecutor) error
	UpdateError() error
}

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

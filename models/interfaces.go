package models

import "github.com/thermokarst/bactdb/Godeps/_workspace/src/github.com/jmoiron/modl"

type updater interface {
	PreUpdate(modl.SqlExecutor) error
	UpdateError() error
}

func Update(u updater) error {
	count, err := DBH.Update(u)
	if err != nil {
		return err
	}
	if count != 1 {
		return u.UpdateError()
	}
	return nil
}

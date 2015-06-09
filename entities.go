package main

import "net/url"

type entity interface {
	marshal() ([]byte, error)
}

type getter interface {
	get(int64, string) (entity, error)
}

type lister interface {
	list(*url.Values) (entity, error)
}

type updater interface {
	update(int64, *entity, Claims) error
	unmarshal([]byte) (entity, error)
}

type creater interface {
	create(*entity, Claims) error
	unmarshal([]byte) (entity, error)
}

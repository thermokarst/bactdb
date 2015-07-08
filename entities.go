package main

import "net/url"

type entity interface {
	marshal() ([]byte, error)
}

type getter interface {
	get(int64, string, Claims) (entity, *appError)
}

type lister interface {
	list(*url.Values, Claims) (entity, *appError)
}

type updater interface {
	update(int64, *entity, string, Claims) *appError
	unmarshal([]byte) (entity, error)
}

type creater interface {
	create(*entity, string, Claims) *appError
	unmarshal([]byte) (entity, error)
}

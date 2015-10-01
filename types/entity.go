package types

type Entity interface {
	Marshal() ([]byte, error)
}

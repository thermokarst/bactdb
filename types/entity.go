package types

// Entity is a a payload or model.
type Entity interface {
	Marshal() ([]byte, error)
}

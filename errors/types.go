package errors

import "errors"

var (
	// ErrSourceNotByteSlice when not a byte-slice.
	ErrSourceNotByteSlice = errors.New("Scan source was not []byte")
)

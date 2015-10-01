package errors

import "errors"

var (
	SourceNotByteSlice = errors.New("Scan source was not []byte")
)

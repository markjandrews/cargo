package proxy

import (
	"errors"
)

// Public Errors
var (
	ErrNotImplemented = errors.New("not implemented")
	ErrInvalidArg = errors.New("invalid argument")
	ErrInvalidData = errors.New("invaid data")
)

func ReverseStringList(values []string){
	if values == nil {
		return
	}

	for i, j := 0, len(values) - 1; i < j; i, j = i+1, j-1 {
		values[i], values[j] = values[j], values[i]
	}
}

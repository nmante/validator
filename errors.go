package main

import (
	"fmt"
	"reflect"
)

// ErrTypeMismatch is returned when the desired type for a validator.Func doesn't match the actual type required
type ErrTypeMismatch struct {
	Value       interface{}
	DesiredType string
}

// Error tells the user the desired type as well as the actual type of the object passed in to a validator
func (e ErrTypeMismatch) Error() string {
	return fmt.Sprintf("%v is of type %s, not %s", e.Value, reflect.TypeOf(e.Value).Name(), e.DesiredType)
}

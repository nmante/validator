package validator

import (
	"errors"
)

var (
	ErrNilValidator = errors.New("Validator must not be nil")
)

type Option func(*Validator) error

func ParallelOption(isParallel bool) Option {
	return func(v *Validator) error {
		if v == nil {
			return ErrNilValidator
		}

		v.enableParallel = true

		return nil
	}
}

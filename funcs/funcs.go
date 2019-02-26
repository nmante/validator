package funcs

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/nmante/validator/compare"
	"github.com/nmante/validator/transform"
	"github.com/nmante/validator/types"
)

var (
	ErrInvalidKind = errors.New("Can't call 'len' on this value")
	validKinds     = map[reflect.Kind]int{
		reflect.Array:  1,
		reflect.Chan:   1,
		reflect.Map:    1,
		reflect.Slice:  1,
		reflect.String: 1,
	}
)

// Response contains info around if a validator function was valid. If it isn't valid, an
// error message is also returned
type Response struct {
	IsValid bool
	Error   string
}

type Interface interface {
	Validate() (Response, error)
}

// Func is function type that all validator functions must follow
type Func func(interface{}) (Response, error)

// IsTransformableTo checks if a value of type 'A' is transformable to type 'B'
func IsTransformableTo(transformer transform.Interface, _type reflect.Type) Func {
	return func(v interface{}) (Response, error) {
		t, err := transformer.Transform(v)
		if err != nil {
			return Response{}, err
		}

		if reflect.TypeOf(t) != _type {
			return Response{
				IsValid: false,
				Error:   fmt.Sprintf("%v not transformable to %s", v, _type.Name()),
			}, nil
		}

		return Response{IsValid: true, Error: ""}, nil
	}
}

// IsEqual transforms a 'v' to a type, and checks if it's equal to 'right'
func IsEqual(transformer transform.Interface, comparer compare.Interface, right interface{}) Func {
	return func(v interface{}) (Response, error) {
		value, err := transformer.Transform(v)
		if err != nil {
			return Response{}, err
		}

		if ok, err := isTypesEqual(value, right); !ok {
			return Response{}, err
		}

		if comparer.Compare(value, right) == 0 {
			return Response{IsValid: true, Error: ""}, nil
		}

		return Response{IsValid: false, Error: fmt.Sprintf("must be equal to %d", right)}, nil
	}
}

// IsBetween checks if a value is between a lower and an upper value
func IsBetween(transformer transform.Interface, comparer compare.Interface, lower interface{}, upper interface{}) Func {
	return func(v interface{}) (Response, error) {
		value, err := transformer.Transform(v)
		if err != nil {
			return Response{}, err
		}

		if ok, err := isTypesEqual(value, lower, upper); !ok {
			return Response{}, err
		}

		if comparer.Compare(lower, value) < 1 && comparer.Compare(value, upper) > -1 {
			return Response{IsValid: true, Error: ""}, nil
		}

		return Response{IsValid: false, Error: fmt.Sprintf("must be between %d and %d", lower, upper)}, nil
	}
}

// IsLength checks if the length of an item equals a value
func IsLength(length int) Func {
	return func(v interface{}) (Response, error) {
		value := reflect.ValueOf(v)
		if _, ok := validKinds[value.Kind()]; !ok {
			return Response{}, ErrInvalidKind
		}

		if length == value.Len() {
			return Response{IsValid: true}, nil
		}

		return Response{
			IsValid: false,
			Error:   fmt.Sprintf("Must have length %d", length),
		}, nil
	}
}

// IsLengthBetween checks if the length of an item is within a range
func IsLengthBetween(lower int, upper int) Func {
	return func(v interface{}) (Response, error) {
		value := reflect.ValueOf(v)
		if _, ok := validKinds[value.Kind()]; !ok {
			return Response{}, ErrInvalidKind
		}

		if lower <= value.Len() && value.Len() <= upper {
			return Response{IsValid: true}, nil
		}

		return Response{
			IsValid: false,
			Error:   fmt.Sprintf("Must be between length %d and %d", lower, upper),
		}, nil
	}
}

// isTypesEqual checks if a list of values all have the same type
func isTypesEqual(values ...interface{}) (bool, error) {
	for i := range values {
		if i == len(values)-1 {
			break
		}

		if reflect.TypeOf(values[i]) != reflect.TypeOf(values[i+1]) {
			return false, errors.New("Values must all have the same type")
		}
	}

	return true, nil
}

// IsType checks if a value is of a certain type
func IsType(_type reflect.Type) Func {
	return func(v interface{}) (Response, error) {
		if reflect.TypeOf(v) != _type {
			return Response{IsValid: false, Error: fmt.Sprintf("must be a %s", _type.Name())}, nil
		}

		return Response{IsValid: true, Error: ""}, nil
	}
}

// IsBool checks if a value is a boolean
func IsBool(v interface{}) (Response, error) {
	return IsType(types.Bool)(v)
}

// IsInt checks if a value is an integer
func IsInt(v interface{}) (Response, error) {
	return IsType(types.Int)(v)
}

// IsInt8 checks if a value is an 8 bit integer
func IsInt8(v interface{}) (Response, error) {
	return IsType(types.Int8)(v)
}

// IsInt16 checks if a value is a 16 bit integer
func IsInt16(v interface{}) (Response, error) {
	return IsType(types.Int16)(v)
}

// IsInt32 checks if a value is a 32 bit integer
func IsInt32(v interface{}) (Response, error) {
	return IsType(types.Int32)(v)
}

// IsInt64 checks if a value is a 64 bit integer
func IsInt64(v interface{}) (Response, error) {
	return IsType(types.Int64)(v)
}

// IsUint checks if a value is an unsigned integer
func IsUint(v interface{}) (Response, error) {
	return IsType(types.Uint)(v)
}

// IsUint8 checks if a value is an 8 bit unsigned integer
func IsUint8(v interface{}) (Response, error) {
	return IsType(types.Uint8)(v)
}

// IsUint16 checks if a value is a 16 bit unsigned integer
func IsUint16(v interface{}) (Response, error) {
	return IsType(types.Uint16)(v)
}

// IsUint32 checks if a value is a 32 bit unsigned integer
func IsUint32(v interface{}) (Response, error) {
	return IsType(types.Uint32)(v)
}

// IsUint64 checks if a value is a 64 bit unsigned integer
func IsUint64(v interface{}) (Response, error) {
	return IsType(types.Uint64)(v)
}

// IsUintptr checks if a value is a unsigned integer pointer
func IsUintptr(v interface{}) (Response, error) {
	return IsType(types.Uintptr)(v)
}

// IsByte checks if a value is a byte
func IsByte(v interface{}) (Response, error) {
	return IsType(types.Byte)(v)
}

// IsRune checks if a value is a rune
func IsRune(v interface{}) (Response, error) {
	return IsType(types.Rune)(v)
}

// IsFloat32 checks if a value is a 32 bit float
func IsFloat32(v interface{}) (Response, error) {
	return IsType(types.Float32)(v)
}

// IsFloat64 checks if a value is a 64 bit float
func IsFloat64(v interface{}) (Response, error) {
	return IsType(types.Float64)(v)
}

// IsComplex64 checks if a value is a 64 bit complex number
func IsComplex64(v interface{}) (Response, error) {
	return IsType(types.Complex64)(v)
}

// IsComplex128 checks if a value is a 128 bit complex number
func IsComplex128(v interface{}) (Response, error) {
	return IsType(types.Complex128)(v)
}

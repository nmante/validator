package types

import (
	"reflect"
)

var _int int
var _int8 int8
var _int16 int16
var _int32 int32
var _int64 int64

var _uint uint
var _uint8 uint8
var _uint16 uint16
var _uint32 uint32
var _uint64 uint64
var _uintptr uintptr

var _byte byte
var _rune rune

var _float32 float32
var _float64 float64

var _complex64 complex64
var _complex128 complex128

var (
	Int   = reflect.TypeOf(_int)
	Int8  = reflect.TypeOf(_int8)
	Int16 = reflect.TypeOf(_int16)
	Int32 = reflect.TypeOf(_int32)
	Int64 = reflect.TypeOf(_int64)

	Uint    = reflect.TypeOf(_uint)
	Uint8   = reflect.TypeOf(_uint8)
	Uint16  = reflect.TypeOf(_uint16)
	Uint32  = reflect.TypeOf(_uint32)
	Uint64  = reflect.TypeOf(_uint64)
	Uintptr = reflect.TypeOf(_uintptr)

	Byte = reflect.TypeOf(_byte)
	Rune = reflect.TypeOf(_rune)

	Float32 = reflect.TypeOf(_float32)
	Float64 = reflect.TypeOf(_float64)

	Complex64  = reflect.TypeOf(_complex64)
	Complex128 = reflect.TypeOf(_complex128)
)

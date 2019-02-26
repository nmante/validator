package funcs

import (
	"github.com/nmante/validator/compare"
	"github.com/nmante/validator/transform"
	"github.com/nmante/validator/types"

	"math/cmplx"
	"strconv"
	"testing"
)

func TestTypeFuncs(t *testing.T) {
	r, err := IsInt(100)
	if err != nil {
		t.Error(err)
	}

	if !r.IsValid {
		t.Error(r.Error)
	}

	r, err = IsComplex128(cmplx.Sqrt(-5 + 12i))
	if err != nil {
		t.Error(err)
	}

	if !r.IsValid {
		t.Error(r.Error)
	}
}

func TestIsTransformableToInt(t *testing.T) {
	r, err := IsTransformableTo(transform.StringToInt, types.Int)("53")
	if err != nil {
		t.Error(err)
	}

	if !r.IsValid {
		t.Error(r.Error)
	}
}

func TestIsEmail(t *testing.T) {
	r, err := String.IsEmail("nii.mante@buzzfeed.com")
	if err != nil {
		t.Error(err)
	}

	if !r.IsValid {
		t.Error(r.Error)
	}

	r, err = String.IsEmail("hello")
	if err != nil {
		t.Error(err)
	}

	if r.IsValid {
		t.Error("'hello' is not an email. This should not be valid")
	}

	r, err = String.IsEmail(3)
	if err == nil {
		t.Error("There should be an error")
	}
}

func TestIsLength(t *testing.T) {
	r, err := IsLength(3)("abc")
	if err != nil {
		t.Error(err)
	}

	if !r.IsValid {
		t.Error(r.Error)
	}

	r, err = IsLength(3)("ab")
	if err != nil {
		t.Error(err)
	}

	if r.IsValid {
		t.Error(r.Error)
	}

	r, err = IsLength(3)(1)
	if err == nil {
		t.Error("There should be an error calling 'len' on invalid kind")
	}
}

func TestIsEqual(t *testing.T) {
	r, err := IsEqual(transform.StringToInt, compare.Int, 100)("53")
	if err != nil {
		t.Error(err)
	}

	if r.IsValid {
		t.Error(r.Error)
	}

	r, err = IsEqual(transform.StringToInt, compare.Int, 100)("100")
	if err != nil {
		t.Error(err)
	}

	if !r.IsValid {
		t.Error(r.Error)
	}
}

func TestIsBetween(t *testing.T) {
	r, err := IsBetween(transform.StringToInt, compare.Int, 1, 100)("53")

	if err != nil {
		t.Error(err)
	}

	if r.IsValid {
		t.Error(r.Error)
	}

	r, err = IsBetween(transform.StringToInt, compare.Int, 1, 100)("101")

	if err != nil {
		t.Error(err)
	}

	if !r.IsValid {
		t.Error(r.Error)
	}
}

func BenchmarkIsBetween(b *testing.B) {
	for k := 0; k < b.N; k++ {
		_, _ = IsBetween(transform.StringToInt, compare.Int, 1, 100)("101")
	}
}

func BenchmarkIsStringBetweenInts(b *testing.B) {
	for k := 0; k < b.N; k++ {
		_, _ = IsStringBetweenInts{value: "101", lower: 1, upper: 100}.Validate()
	}
}

func BenchmarkNormalIsBetween(b *testing.B) {
	isBetween := func(value string, lower, upper int) (bool, error) {
		parsed, err := strconv.ParseInt(value, 10, 0)
		if err != nil {
			return false, err
		}

		if lower <= int(parsed) && int(parsed) <= upper {
			return true, nil
		}

		return false, nil
	}

	for k := 0; k < b.N; k++ {
		_, _ = isBetween("101", 1, 100)
	}
}

package validator

import (
	"math/cmplx"
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
	r, err := IsTransformableToInt(StringToInt{})("53")
	if err != nil {
		t.Error(err)
	}

	if !r.IsValid {
		t.Error(r.Error)
	}
}

func TestIsEmail(t *testing.T) {
	r, err := IsEmail("nii.mante@buzzfeed.com")
	if err != nil {
		t.Error(err)
	}

	if !r.IsValid {
		t.Error(r.Error)
	}

	r, err = IsEmail("hello")
	if err != nil {
		t.Error(err)
	}

	if r.IsValid {
		t.Error("'hello' is not an email. This should not be valid")
	}

	r, err = IsEmail(3)
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
	r, err := IsEqual(StringToInt{}, IntComparer{}, 100)("53")
	if err != nil {
		t.Error(err)
	}

	if r.IsValid {
		t.Error(r.Error)
	}

	r, err = IsEqual(StringToInt{}, IntComparer{}, 100)("100")
	if err != nil {
		t.Error(err)
	}

	if !r.IsValid {
		t.Error(r.Error)
	}
}

func TestIsBetween(t *testing.T) {
	r, err := IsBetween(StringToInt{}, IntComparer{}, 1, 100)("53")

	if err != nil {
		t.Error(err)
	}

	if r.IsValid {
		t.Error(r.Error)
	}

	r, err = IsBetween(StringToInt{}, IntComparer{}, 1, 100)("101")

	if err != nil {
		t.Error(err)
	}

	if !r.IsValid {
		t.Error(r.Error)
	}
}

package transform

import (
	"testing"
)

func TestStringToTranformers(t *testing.T) {
	stringTranformerTests := []struct {
		val           string
		transformer   Interface
		transformed   interface{}
		shouldBeError bool
	}{
		{val: "1", transformer: StringToInt, transformed: 1, shouldBeError: false},
		{val: "4", transformer: StringToInt, transformed: 4, shouldBeError: false},
		{val: "a", transformer: StringToInt, transformed: nil, shouldBeError: true},
		{val: "4.5", transformer: StringToFloat32, transformed: 4.5, shouldBeError: false},
		{val: "10000000000000000.5", transformer: StringToFloat64, transformed: 10000000000000000.5, shouldBeError: false},
		{val: "true", transformer: StringToBool, transformed: true, shouldBeError: false},
		{val: "8", transformer: StringToUint, transformed: uint64(8), shouldBeError: false},
	}

	for _, test := range stringTranformerTests {
		val, err := test.transformer.Transform(test.val)

		if test.shouldBeError && err == nil {
			t.Errorf("There should be an error")
		}

		if !test.shouldBeError && val != test.transformed {
			t.Errorf("Transformed value %s does not equal test value %s", val, test.transformed)
		}
	}
}

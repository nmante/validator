package compare

import (
	"testing"
)

func TestComparers(t *testing.T) {
	intCompareTests := []struct {
		left     interface{}
		right    interface{}
		comparer Interface
		result   interface{}
	}{
		{left: 1, right: 1, comparer: Int, result: 0},
		{left: 1, right: 2, comparer: Int, result: -1},
		{left: float32(4.5), right: float32(3.55), comparer: Float32, result: 1},
		{left: float32(1.1), right: float32(1.1), comparer: Float32, result: 0},
		{left: 1.3, right: 1.5, comparer: Float64, result: -1},
	}

	for _, test := range intCompareTests {
		if result := test.comparer.Compare(test.left, test.right); result != test.result {
			t.Errorf("Compare result is: %d. Actual result should be %d", result, test.result)
		}
	}
}

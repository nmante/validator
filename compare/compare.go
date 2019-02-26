package compare

var (
	Default = _default{}
	Int     = _int{}
	Float32 = _float32{}
	Float64 = _float64{}
)

// Comparer compares left & right and returns -1 (less than), 0 (equal), 1 (greater than)
type Interface interface {
	Compare(left interface{}, right interface{}) int
}

type _default struct{}

func (i _default) Compare(left interface{}, right interface{}) int {
	l := left.(int)
	r := right.(int)

	if l < r {
		return -1
	} else if l == r {
		return 0
	}

	return 1
}

type _int struct{}

func (i _int) Compare(left interface{}, right interface{}) int {
	l := left.(int)
	r := right.(int)

	if l < r {
		return -1
	} else if l == r {
		return 0
	}

	return 1
}

type _float32 struct{}

func (f _float32) Compare(left interface{}, right interface{}) int {
	l := left.(float32)
	r := right.(float32)

	if l < r {
		return -1
	} else if l == r {
		return 0
	}

	return 1
}

type _float64 struct{}

func (f _float64) Compare(left interface{}, right interface{}) int {
	l := left.(float64)
	r := right.(float64)

	if l < r {
		return -1
	} else if l == r {
		return 0
	}

	return 1
}

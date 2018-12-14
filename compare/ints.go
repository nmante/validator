package compare

var (
	Int = _int{}
)

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

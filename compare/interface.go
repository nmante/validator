package compare

// Comparer compares left & right and returns -1, 0, 1 for less than, equal, greater than
type Interface interface {
	Compare(left interface{}, right interface{}) int
}

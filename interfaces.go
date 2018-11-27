package validator

// Transformer allows us to convert a value 'A' to a different value 'B'.
// A and B's types/values can be the same or differ
type Transformer interface {
	Transform(v interface{}) (interface{}, error)
}

// Comparer compares left & right and returns -1, 0, 1 for less than, equal, greater than
type Comparer interface {
	Compare(left interface{}, right interface{}) int
}

package transform

// Transformer allows us to convert a value 'A' to a different value 'B'.
// A and B's types/values can be the same or differ
type Interface interface {
	Transform(v interface{}) (interface{}, error)
}

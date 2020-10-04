package datastructures

type (
	ItemType interface{}

	KeyValuePair struct {
		Key   interface{}
		Value interface{}
	}

	IndexValuePair struct {
		Index int
		Value interface{}
	}

	IndexIterator interface {
		Iterate() <-chan IndexValuePair
	}

	KeyIterator interface {
		IterateKeys() <-chan KeyValuePair
	}

	Joiner interface {
		Join(separator string) string
	}

	Sizer interface {
		Len() int
		Cap() int
	}

	Emptier interface {
		Empty()
		IsEmpty() bool
	}
)

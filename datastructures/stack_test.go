package datastructures

import "testing"

func TestStack_Iterator(t *testing.T) {
	s := NewStack("a", "b", "c", "d")

	for elem := range s.Iterate() {
		t.Log(elem.Index, elem.Value)
	}

	t.Log("END")
}

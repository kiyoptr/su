package reflect

import (
	"testing"
)

func TestReverseSlice(t *testing.T) {
	s1 := []int{1, 2, 3}
	s2, err := ReverseSlice(s1)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(s1)
	t.Log(s2)
}

package errors

import (
	"io"
	"testing"
)

func TestHasError(t *testing.T) {
	var err error = Newi(io.EOF, "Inner most error")
	err = Newi(err, "Some inner error")
	err = Newi(err, "First error")

	if !Is(err, io.EOF) {
		t.Error("err doesn't contain io.EOF")
	}

	t.Log(err)
}

package errors

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

type ErrTest int

func (e ErrTest) Error() string {
	return fmt.Sprintf("%d", e)
}

func TestHasError(t *testing.T) {
	err := Newi(
		Newi(
			Newi(
				//Newi(&net.OpError{Op: "test", Net: "tcp", Err: io.EOF}, "EOF error"),
				Newi(ErrTest(150), "Err Test"),
				"one more to last"),
			"two more to last"),
		"three more to last")

	if !Is(err, io.EOF) {
		t.Error("stdlib Is failed")
	}

	if err := ioutil.WriteFile("log.txt", []byte(err.Error()), os.ModePerm); err != nil {
		t.Error(err)
	}
}

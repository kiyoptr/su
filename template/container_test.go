package template

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// To run this test, working directory must be set to directory containing this test file.
func TestContainer_LoadDirectory(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	container := NewContainer(nil)
	if err := container.LoadDirectory(filepath.Join(wd, "templates")); err != nil {
		t.Fatal(err)
	}

	s := strings.Builder{}
	if err := container.Execute(&s, "inner.test", nil); err != nil {
		t.Fatal(err)
	}

	t.Log(s.String())
}

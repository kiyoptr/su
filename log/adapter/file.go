package adapter

import (
	"fmt"
	"os"
)

type file struct {
	output     *os.File
	customFile bool
}

func File(name string) *file {
	f, err := os.Create(name)
	if err != nil {
		return nil
	}

	return &file{
		output:     f,
		customFile: true,
	}
}

func Stdout() *file {
	return &file{output: os.Stdout}
}

func Stderr() *file {
	return &file{output: os.Stderr}
}

func (f *file) Write(message string) {
	fmt.Fprintln(f.output, message)
}

func (f *file) Close() error {
	if f.customFile {
		return f.output.Close()
	}

	return nil
}

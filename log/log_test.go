package log

import (
	"github.com/kiyoptr/su/log/adapter"
	"github.com/kiyoptr/su/log/tagprovider"
	"sync"
	"testing"
)

func TestNew(t *testing.T) {
	l, err := New().
		Name("test logger").
		WithAdapters(adapter.Stderr()).
		WithDefaultMode(Error).
		WithTags(tagprovider.DateTime("Mon Jan 2 15:04:05 -0700 MST 2006")).
		Build()

	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	l.Mode(Info).Tag("chos", "bus").Write()
	l.Writef("ass")
}

func TestThreadSafety(t *testing.T) {
	var err error
	instance, err = New().
		Name("tts").
		WithAdapters(adapter.Stderr()).
		WithDefaultMode(Info).
		WithTags(tagprovider.DateTime("2006/Jan/2 15:04:05 -0700")).
		Build()

	if err != nil {
		t.Fatal(err)
	}
	defer instance.Close()

	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			instance.Mode(Warning).Writef("this is from thread %d", i)
		}(i)
	}
	wg.Wait()
}

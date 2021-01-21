package log

import (
	"errors"
	"github.com/kiyoptr/su/log/adapter"
	"github.com/kiyoptr/su/log/tagprovider"
	"sync"
)

type Builder struct {
	adapters []adapter.Adapter
	tl       taglist
}

func New() *Builder {
	return &Builder{
		tl: newTagList(),
	}
}

func (b *Builder) Name(name string) *Builder {
	b.tl.set(toName, tagprovider.Constant("name", name))
	return b
}

func (b *Builder) WithDefaultMode(mode Mode) *Builder {
	b.tl.set(toMode, tagprovider.Constant("mode", mode))
	return b
}

func (b *Builder) WithAdapters(a ...adapter.Adapter) *Builder {
	if b.adapters == nil {
		b.adapters = make([]adapter.Adapter, 0, len(a))
	}

	b.adapters = append(b.adapters, a...)

	return b
}

func (b *Builder) WithTags(providers ...tagprovider.Provider) *Builder {
	for _, p := range providers {
		b.tl.set(toCustom, p)
	}

	return b
}

var (
	ErrNoAdaptersProvided = errors.New("No adapters are provided")
)

func (b *Builder) Build() (l *Logger, err error) {
	if len(b.adapters) == 0 {
		err = ErrNoAdaptersProvided
		return
	}

	l = &Logger{
		staticTags: b.tl,
		customTags: newTagList(),
		adapters:   b.adapters,
		lock:       sync.Mutex{},
	}

	return
}

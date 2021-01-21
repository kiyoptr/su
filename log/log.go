package log

import (
	"fmt"
	"github.com/kiyoptr/su/log/adapter"
	"github.com/kiyoptr/su/log/tagprovider"
	"sync"
)

var (
	// instance is the global instance of logger
	instance *Logger
)

func SetGlobal(l *Logger) {
	instance = l
}

func Instance() *Logger { return instance }

// Logger is a thread-safe logging type
type Logger struct {
	staticTags taglist
	customTags taglist
	adapters   []adapter.Adapter
	lock       sync.Mutex
}

func (l *Logger) Close() error {
	for _, a := range l.adapters {
		if err := a.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (l *Logger) Mode(mode Mode) *Logger {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.customTags.set(toMode, tagprovider.Constant("mode", mode))
	return l
}

func (l *Logger) Tag(key string, value interface{}) *Logger {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.customTags.set(toCustom, tagprovider.Constant(key, value))
	return l
}

func (l *Logger) Write() {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.write()
}

func (l *Logger) Writef(format string, args ...interface{}) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.customTags.set(toMessage, tagprovider.Constant("message", fmt.Sprintf(format, args...)))
	l.write()
}

func (l *Logger) write() {
	t := l.staticTags.merge(l.customTags).build()
	for _, a := range l.adapters {
		a.Write(t)
	}
	l.customTags = newTagList()
}

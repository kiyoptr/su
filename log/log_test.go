package log

import (
	"fmt"
	"sync"
	"testing"
)

func TestLogHard(t *testing.T) {
	Load("test_logs")
	wg := sync.WaitGroup{}
	max := 1000
	wg.Add(max)
	for i := 0; i < max; i++ {
		go func(n int) {
			Log(fmt.Sprintf("LOG MSG %v", n))
			wg.Done()
		}(i)
	}
	wg.Wait()
	Close()
}

func TestLogThreadSafety(t *testing.T) {

}

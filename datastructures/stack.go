package datastructures

import (
	"fmt"
	"strings"
	"sync"
)

var StackInitialSize = 10

type Stack struct {
	items []ItemType
	lock  sync.RWMutex
}

func NewStack(items ...ItemType) *Stack {
	return &Stack{items, sync.RWMutex{}}
}

func (s *Stack) Push(value interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.items == nil {
		s.items = make([]ItemType, StackInitialSize)
	}

	s.items = append(s.items, value)
}

func (s *Stack) Pop() interface{} {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.items == nil || len(s.items) == 0 {
		return nil
	}

	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]

	return item
}

func (s *Stack) Peek() interface{} {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if s.items == nil || len(s.items) == 0 {
		return nil
	}

	return s.items[len(s.items)-1]
}

func (s *Stack) Len() int {
	if s.items == nil {
		return 0
	}

	s.lock.RLock()
	defer s.lock.RUnlock()

	return len(s.items)
}

func (s *Stack) Cap() int {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return cap(s.items)
}

func (s *Stack) Empty() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.items = []ItemType{}
}

func (s *Stack) IsEmpty() bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.items == nil || len(s.items) == 0
}

func (s *Stack) Iterate() <-chan IndexValuePair {
	it := make(chan IndexValuePair, 1)

	go func() {
		s.lock.RLock()
		defer s.lock.RUnlock()

		for i, v := range s.items {
			it <- IndexValuePair{i, v}
		}

		close(it)
	}()

	return it
}

func (s *Stack) Join(sep string) string {
	values := make([]string, s.Len())

	for elem := range s.Iterate() {
		values[elem.Index] = fmt.Sprintf("%v", elem.Value)
	}

	return strings.Join(values, sep)
}

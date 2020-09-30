package datastructures

import "sync"

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
		s.items = []ItemType{}
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

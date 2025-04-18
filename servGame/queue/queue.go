package queue

import (
	"errors"
	"net"
	"sync"
)

type Queue struct {
	lock sync.Mutex // you don't have to do this if you don't want thread safety
	s    []net.Conn
}

func NewQueue() *Queue {
	return &Queue{sync.Mutex{}, make([]net.Conn, 0)}
}

func (s *Queue) Add(v net.Conn) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.s = append(s.s, v)
}

func (s *Queue) Remove() (net.Conn, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	l := len(s.s)
	if l == 0 {
		return nil, errors.New("empty stack")
	}

	res := s.s[0]
	s.s = s.s[1:l]
	return res, nil
}

func (s *Queue) Count() int {
	s.lock.Lock()
	defer s.lock.Unlock()

	return len(s.s)
}

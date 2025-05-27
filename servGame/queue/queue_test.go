package queue_test

import (
	"SERV/queue"
	"net"
	"testing"
	"time"
)

// MockConn это mock-реализация net.Conn для тестирования
type MockConn struct{}

func (m *MockConn) Read(b []byte) (n int, err error)   { return 0, nil }
func (m *MockConn) Write(b []byte) (n int, err error)  { return 0, nil }
func (m *MockConn) Close() error                       { return nil }
func (m *MockConn) LocalAddr() net.Addr                { return nil }
func (m *MockConn) RemoteAddr() net.Addr               { return nil }
func (m *MockConn) SetDeadline(t time.Time) error      { return nil }
func (m *MockConn) SetReadDeadline(t time.Time) error  { return nil }
func (m *MockConn) SetWriteDeadline(t time.Time) error { return nil }

func TestNewQueue(t *testing.T) {
	q := queue.NewQueue()
	if q == nil {
		t.Error("NewQueue() returned nil")
	}
	if q.Count() != 0 {
		t.Errorf("New queue should be empty, got %d", q.Count())
	}
}

func TestQueue_Add(t *testing.T) {
	q := queue.NewQueue()
	conn := &MockConn{}

	q.Add(conn)
	if q.Count() != 1 {
		t.Errorf("Expected count 1 after Add, got %d", q.Count())
	}
}

func TestQueue_Remove(t *testing.T) {
	t.Run("remove from non-empty queue", func(t *testing.T) {
		q := queue.NewQueue()
		conn := &MockConn{}
		q.Add(conn)

		removed, err := q.Remove()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if removed != conn {
			t.Error("Removed connection doesn't match added one")
		}
		if q.Count() != 0 {
			t.Error("Queue should be empty after remove")
		}
	})

	t.Run("remove from empty queue", func(t *testing.T) {
		q := queue.NewQueue()

		_, err := q.Remove()
		if err == nil {
			t.Error("Expected error when removing from empty queue")
		}
		if err.Error() != "empty stack" {
			t.Errorf("Unexpected error message: %v", err)
		}
	})
}

func TestQueue_Count(t *testing.T) {
	q := queue.NewQueue()
	if q.Count() != 0 {
		t.Error("New queue should have count 0")
	}

	conn1 := &MockConn{}
	conn2 := &MockConn{}

	q.Add(conn1)
	if q.Count() != 1 {
		t.Error("Count should be 1 after first add")
	}

	q.Add(conn2)
	if q.Count() != 2 {
		t.Error("Count should be 2 after second add")
	}

	q.Remove()
	if q.Count() != 1 {
		t.Error("Count should be 1 after remove")
	}
}

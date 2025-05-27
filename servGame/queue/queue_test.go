package queue_test

import (
	"SERV/queue"
	"testing"
)

func TestQueueCreate(t *testing.T) {
	q := queue.NewQueue()

	if q.Count() != 0 {
		t.Fail()
	}
}

func TestQueueAdd(t *testing.T) {

}

package tests

import (
	"testing"

	"github.com/google/uuid"
	"github.com/kokaq/core/queue"
	"github.com/stretchr/testify/assert"
)

func TestQueue_EnqueueDequeue(t *testing.T) {
	q, _ := queue.NewDefaultKokaq(1, 1)
	assert.True(t, q.IsEmpty(), "Queue should be empty initially")

	a := uuid.New()
	b := uuid.New()
	c := uuid.New()

	q.PushItem(queue.NewKokaqItem(c, 1))
	q.PushItem(queue.NewKokaqItem(b, 2))
	q.PushItem(queue.NewKokaqItem(a, 3))

	assert.False(t, q.IsEmpty(), "Queue should not be empty after enqueuing")

	item, err := q.PopItem()
	assert.NotNil(t, err, "Dequeue should succeed")
	assert.Equal(t, a, item.Id, "First dequeued item should be highest priority 'a'")

	item, err = q.PopItem()
	assert.NotNil(t, err, "Dequeue should succeed")
	assert.Equal(t, b, item.Id, "Second dequeued item should be highest priority 'b'")

	item, err = q.PopItem()
	assert.NotNil(t, err, "Dequeue should succeed")
	assert.Equal(t, c, item.Id, "Third dequeued item should be highest priority 'c'")

	assert.True(t, q.IsEmpty(), "Queue should be empty after dequeuing all items")
}

func TestQueue_DequeueEmpty(t *testing.T) {
	q, _ := queue.NewDefaultKokaq(2, 2)
	item, err := q.PopItem()
	assert.NotNil(t, err, "Dequeue should fail on empty queue")
	assert.Nil(t, item, "Dequeued item from empty queue should be nil")
}

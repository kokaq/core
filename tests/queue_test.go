package tests

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kokaq/core/queue"
	"github.com/stretchr/testify/assert"
)

func setupTestQueue(t *testing.T, enableDLQ, enableInvisible bool) (*queue.Queue, func()) {
	tmpDir := filepath.Join(os.TempDir(), "testqueue", uuid.New().String())
	config := queue.QueueConfiguration{
		QueueName:       "test",
		QueueId:         1,
		EnableDLQ:       enableDLQ,
		EnableInvisible: enableInvisible,
	}
	q, err := queue.NewQueue(tmpDir, config)
	assert.NoError(t, err)
	cleanup := func() {
		_ = q.Delete()
		_ = os.RemoveAll(tmpDir)
	}
	return q, cleanup
}

func TestNewQueueAndDelete(t *testing.T) {
	q, cleanup := setupTestQueue(t, true, true)
	defer cleanup()
	assert.NotNil(t, q)
	assert.Equal(t, uint32(1), q.Id)
	assert.Equal(t, "test", q.Name)
	assert.True(t, q.EnableDLQ)
	assert.True(t, q.EnableInvisible)
	err := q.Delete()
	assert.NoError(t, err)
}

func TestQueueIsEmpty(t *testing.T) {
	q, cleanup := setupTestQueue(t, false, false)
	defer cleanup()
	empty, err := q.IsEmpty()
	assert.NoError(t, err)
	assert.True(t, empty)
}

func TestQueueEnqueueDequeuePeek(t *testing.T) {
	q, cleanup := setupTestQueue(t, false, false)
	defer cleanup()
	item := &queue.QueueItem{
		MessageId: uuid.New(),
		Priority:  1,
	}
	err := q.Enqueue(item)
	assert.NoError(t, err)

	empty, err := q.IsEmpty()
	assert.NoError(t, err)
	assert.False(t, empty)

	peeked, err := q.Peek()
	assert.NoError(t, err)
	assert.Equal(t, item.MessageId, peeked.MessageId)

	dequeued, err := q.Dequeue()
	assert.NoError(t, err)
	assert.Equal(t, item.MessageId, dequeued.MessageId)

	empty, err = q.IsEmpty()
	assert.NoError(t, err)
	assert.True(t, empty)
}

func TestQueuePeekLockAndAck(t *testing.T) {
	q, cleanup := setupTestQueue(t, false, true)
	defer cleanup()
	item := &queue.QueueItem{
		MessageId: uuid.New(),
		Priority:  2,
	}
	err := q.Enqueue(item)
	assert.NoError(t, err)

	peeked, lockId, err := q.PeekLock()
	assert.NoError(t, err)
	assert.Equal(t, item.MessageId, peeked.MessageId)
	assert.Equal(t, "", lockId)

	err = q.Ack(lockId)
	assert.NoError(t, err)
}

func TestQueueDLQOperations(t *testing.T) {
	q, cleanup := setupTestQueue(t, true, false)
	defer cleanup()
	item := &queue.QueueItem{
		MessageId: uuid.New(),
		Priority:  3,
	}
	err := q.Enqueue(item)
	assert.NoError(t, err)

	err = q.MoveToDLQ(item.MessageId)
	assert.NoError(t, err)

	dlqItems, err := q.PeekDLQ()
	assert.NoError(t, err)
	assert.Nil(t, dlqItems) // stub returns nil

	dequeued, err := q.DequeueDLQ()
	assert.NoError(t, err)
	assert.Nil(t, dequeued) // stub returns nil

	err = q.MoveFromDLQ(item.MessageId)
	assert.NoError(t, err)

	err = q.ClearDLQ()
	assert.NoError(t, err)
}

func TestQueueStatsAndClear(t *testing.T) {
	q, cleanup := setupTestQueue(t, true, true)
	defer cleanup()
	stats, err := q.GetStats()
	assert.NoError(t, err)
	assert.NotNil(t, stats)

	err = q.Clear()
	assert.NoError(t, err)
}

func TestQueueListMessages(t *testing.T) {
	q, cleanup := setupTestQueue(t, false, false)
	defer cleanup()
	items, err := q.ListMessages()
	assert.NoError(t, err)
	assert.Nil(t, items)

	locked, err := q.ListLockedMessages()
	assert.NoError(t, err)
	assert.Nil(t, locked)

	dlq, err := q.ListDLQMessages()
	assert.NoError(t, err)
	assert.Nil(t, dlq)
}

func TestQueueVisibilityTimeoutMethods(t *testing.T) {
	q, cleanup := setupTestQueue(t, false, true)
	defer cleanup()
	lockId := "dummy-lock"
	err := q.Extend(lockId, time.Second)
	assert.NoError(t, err)
	err = q.SetVisibilityTimeout(time.Second)
	assert.NoError(t, err)
	err = q.RefreshVisibilityTimeout(lockId)
	assert.NoError(t, err)
	err = q.ReleaseLock(lockId)
	assert.NoError(t, err)
	expired, err := q.IsExpired(lockId)
	assert.NoError(t, err)
	assert.False(t, expired)
}

func TestQueueAutoMoveToDLQ(t *testing.T) {
	q, cleanup := setupTestQueue(t, true, false)
	defer cleanup()
	messageId := uuid.New()
	err := q.AutoMoveToDLQ(messageId, 5)
	assert.NoError(t, err)
}

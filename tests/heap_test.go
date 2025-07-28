package tests

import (
	"testing"

	"github.com/google/uuid"
	"github.com/kokaq/core/queue"
	"github.com/stretchr/testify/assert"
)

func TestHeapInitialize(t *testing.T) {
	tmpDir := t.TempDir()
	heap, err := queue.NewHeap(tmpDir, 4, 8, 8, 16)
	assert.True(t, heap != nil && err == nil, "heap should be initialized")
}

func TestHeapEnqueue(t *testing.T) {
	tmpDir := t.TempDir()
	heap, err := queue.NewHeap(tmpDir, 4, 8, 8, 16)
	if err != nil {
		t.Fatalf("Failed to initialize heap: %v", err)
	}

	msgID := uuid.New()
	item := &queue.QueueItem{
		MessageId: msgID,
		Priority:  10,
	}

	// Enqueue
	if err := heap.Enqueue(item); err != nil {
		t.Fatalf("Enqueue failed: %v", err)
	}
}

func TestHeapDequeue(t *testing.T) {
	tmpDir := t.TempDir()
	heap, err := queue.NewHeap(tmpDir, 4, 8, 8, 16)
	if err != nil {
		t.Fatalf("Failed to initialize heap: %v", err)
	}

	msgID := uuid.New()
	item := &queue.QueueItem{
		MessageId: msgID,
		Priority:  10,
	}

	// Enqueue
	if err := heap.Enqueue(item); err != nil {
		t.Fatalf("Enqueue failed: %v", err)
	}

	// Dequeue
	dequeued, err := heap.Dequeue()
	if err != nil {
		t.Fatalf("Dequeue failed: %v", err)
	}
	if dequeued.MessageId != msgID || dequeued.Priority != 10 {
		t.Errorf("Dequeue returned wrong item: got %+v", dequeued)
	}
}

func TestHeapPeek(t *testing.T) {
	tmpDir := t.TempDir()
	heap, err := queue.NewHeap(tmpDir, 4, 8, 8, 16)
	if err != nil {
		t.Fatalf("Failed to initialize heap: %v", err)
	}

	msgID := uuid.New()
	item := &queue.QueueItem{
		MessageId: msgID,
		Priority:  10,
	}

	// Enqueue
	if err := heap.Enqueue(item); err != nil {
		t.Fatalf("Enqueue failed: %v", err)
	}

	// Peek
	peeked, err := heap.Peek()
	if err != nil {
		t.Fatalf("Peek failed: %v", err)
	}
	if peeked.MessageId != msgID || peeked.Priority != 10 {
		t.Errorf("Peek returned wrong item: got %+v", peeked)
	}
}

// func TestHeapIsEmpty(t *testing.T) {
// 	tmpDir := t.TempDir()
// 	heap, err := queue.NewHeap(tmpDir, 4, 8, 8, 16)
// 	if err != nil {
// 		t.Fatalf("Failed to initialize heap: %v", err)
// 	}

// 	empty, err := heap.IsEmpty()
// 	if err != nil {
// 		t.Fatalf("IsEmpty failed: %v", err)
// 	}
// 	if !empty {
// 		t.Error("Heap should be empty after initialization")
// 	}

// 	item := &queue.QueueItem{
// 		MessageId: uuid.New(),
// 		Priority:  5,
// 	}
// 	if err := heap.Enqueue(item); err != nil {
// 		t.Fatalf("Enqueue failed: %v", err)
// 	}

// 	empty, err = heap.IsEmpty()
// 	if err != nil {
// 		t.Fatalf("IsEmpty failed: %v", err)
// 	}
// 	if empty {
// 		t.Error("Heap should not be empty after enqueue")
// 	}
// }

func TestHeapEnqueueZeroPriority(t *testing.T) {
	tmpDir := t.TempDir()
	heap, err := queue.NewHeap(tmpDir, 4, 8, 8, 16)
	if err != nil {
		t.Fatalf("Failed to initialize heap: %v", err)
	}

	item := &queue.QueueItem{
		MessageId: uuid.New(),
		Priority:  0,
	}
	err = heap.Enqueue(item)
	if err == nil {
		t.Error("Expected error when enqueuing zero priority")
	}
}

func TestHeapDequeueEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	heap, err := queue.NewHeap(tmpDir, 4, 8, 8, 16)
	if err != nil {
		t.Fatalf("Failed to initialize heap: %v", err)
	}

	_, err = heap.Dequeue()
	if err == nil {
		t.Error("Expected error when dequeuing from empty heap")
	}
}

func TestHeapPeekEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	heap, err := queue.NewHeap(tmpDir, 4, 8, 8, 16)
	if err != nil {
		t.Fatalf("Failed to initialize heap: %v", err)
	}

	_, err = heap.Peek()
	if err == nil {
		t.Error("Expected error when peeking from empty heap")
	}
}

func TestHeapPersistence(t *testing.T) {
	tmpDir := t.TempDir()
	heap, err := queue.NewHeap(tmpDir, 4, 8, 8, 16)
	if err != nil {
		t.Fatalf("Failed to initialize heap: %v", err)
	}

	msgID := uuid.New()
	item := &queue.QueueItem{
		MessageId: msgID,
		Priority:  7,
	}
	if err := heap.Enqueue(item); err != nil {
		t.Fatalf("Enqueue failed: %v", err)
	}

	// Simulate reload
	heap2, err := queue.NewHeap(tmpDir, 4, 8, 8, 16)
	if err != nil {
		t.Fatalf("Failed to reload heap: %v", err)
	}
	peeked, err := heap2.Peek()
	if err != nil {
		t.Fatalf("Peek after reload failed: %v", err)
	}
	if peeked.MessageId != msgID || peeked.Priority != 7 {
		t.Errorf("Peek after reload returned wrong item: got %+v", peeked)
	}
}

// func TestHeapCleanup(t *testing.T) {
// 	tmpDir := t.TempDir()
// 	heap, err := queue.NewHeap(tmpDir, 4, 8, 8, 16)
// 	if err != nil {
// 		t.Fatalf("Failed to initialize heap: %v", err)
// 	}

// 	item := &queue.QueueItem{
// 		MessageId: uuid.New(),
// 		Priority:  3,
// 	}
// 	if err := heap.Enqueue(item); err != nil {
// 		t.Fatalf("Enqueue failed: %v", err)
// 	}

// 	_, err = heap.Dequeue()
// 	if err != nil {
// 		t.Fatalf("Dequeue failed: %v", err)
// 	}

// 	empty, err := heap.IsEmpty()
// 	if err != nil {
// 		t.Fatalf("IsEmpty failed: %v", err)
// 	}
// 	if !empty {
// 		t.Error("Heap should be empty after dequeue")
// 	}

// 	// Check index file deleted
// 	indexPath := filepath.Join(heap.GetConfig().IndexPath, "3")
// 	if _, err := os.Stat(indexPath); !os.IsNotExist(err) {
// 		t.Error("Index file should be deleted after last dequeue")
// 	}
// }

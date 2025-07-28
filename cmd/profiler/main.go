package main

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/kokaq/core/profiler"
	"github.com/kokaq/core/queue"
)

func main() {
	stop := profiler.Start(profiler.Config{
		CPUProfilePath:    "cpu.prof",
		MemProfilePath:    "mem.prof",
		BlockProfilePath:  "block.prof",
		GoroutineDumpPath: "goroutines.prof",
		TracePath:         "trace.out",
	})
	defer stop()
	ns := queue.NewNamespace("C://code/kokaq/bin", queue.NamespaceConfig{
		NamespaceName: "data-db",
		NamespaceId:   1,
	})

	fmt.Println("Namespace created: #", ns.Id, ": ", ns.Name)

	qConfig := queue.QueueConfiguration{
		QueueName:       "test-queue",
		QueueId:         1,
		EnableDLQ:       true,
		EnableInvisible: true,
	}
	var q *queue.Queue
	var err error

	if q, err = ns.AddQueue(&qConfig); err != nil {
		fmt.Println("Error adding queue:", err)
	} else {
		fmt.Println("Queue added: #", q.Id, ": ", q.Name)
	}

	// if empty, _ := q.IsEmpty(); empty {
	// 	fmt.Println("Queue is Empty")
	// }

	var qi = &queue.QueueItem{
		MessageId: uuid.New(),
		Priority:  1,
	}

	if err = q.Enqueue(qi); err != nil {
		fmt.Println("Error enqueuing item:", err)
	} else {
		fmt.Println("Enqueued item:", qi.MessageId, "with priority", qi.Priority)
	}

	qi = &queue.QueueItem{
		MessageId: uuid.New(),
		Priority:  2,
	}
	if err = q.Enqueue(qi); err != nil {
		fmt.Println("Error enqueuing item:", err)
	} else {
		fmt.Println("Enqueued item:", qi.MessageId, "with priority", qi.Priority)
	}

	qi = &queue.QueueItem{
		MessageId: uuid.New(),
		Priority:  2,
	}
	if err = q.Enqueue(qi); err != nil {
		fmt.Println("Error enqueuing item:", err)
	} else {
		fmt.Println("Enqueued item:", qi.MessageId, "with priority", qi.Priority)
	}

	qi = &queue.QueueItem{
		MessageId: uuid.New(),
		Priority:  3,
	}
	if err = q.Enqueue(qi); err != nil {
		fmt.Println("Error enqueuing item:", err)
	} else {
		fmt.Println("Enqueued item:", qi.MessageId, "with priority", qi.Priority)
	}

	if qi, err := q.Dequeue(); err != nil {
		fmt.Println("Error dequeuing item:", err)
	} else {
		fmt.Println("Dequeued item:", qi.MessageId, "with priority", qi.Priority)
	}

	if qi, err := q.Dequeue(); err != nil {
		fmt.Println("Error dequeuing item:", err)
	} else {
		fmt.Println("Dequeued item:", qi.MessageId, "with priority", qi.Priority)
	}

	if qi, err := q.Peek(); err != nil {
		fmt.Println("Error peeking item:", err)
	} else {
		fmt.Println("Peeked item:", qi.MessageId, "with priority", qi.Priority)
	}
	ns.DeleteQueue(1)
}

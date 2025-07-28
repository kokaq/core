package main

import (
	"github.com/google/uuid"
	"github.com/kokaq/core/queue"
)

func main() {

	ns := queue.NewNamespace("C://code/kokaq/bin", queue.NamespaceConfig{
		NamespaceName: "data-db",
		NamespaceId:   1,
	})
	qConfig := queue.QueueConfiguration{
		QueueName:       "test-queue",
		QueueId:         1,
		EnableDLQ:       true,
		EnableInvisible: true,
	}
	var q *queue.Queue
	q, _ = ns.AddQueue(&qConfig)
	// if empty, _ := q.IsEmpty(); empty {
	// 	fmt.Println("Queue is Empty")
	// }
	var qi = &queue.QueueItem{
		MessageId: uuid.New(),
		Priority:  1,
	}
	_ = q.Enqueue(qi)
	qi = &queue.QueueItem{
		MessageId: uuid.New(),
		Priority:  2,
	}
	_ = q.Enqueue(qi)
	qi = &queue.QueueItem{
		MessageId: uuid.New(),
		Priority:  2,
	}
	_ = q.Enqueue(qi)
	qi = &queue.QueueItem{
		MessageId: uuid.New(),
		Priority:  3,
	}
	_ = q.Enqueue(qi)
	q.Dequeue()
	q.Dequeue()
	q.Peek()
	ns.DeleteQueue(1)
}

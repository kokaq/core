package main

import (
	"github.com/google/uuid"
	"github.com/kokaq/core/internals/logger"
	"github.com/kokaq/core/queue"
)

func main() {

	ns := queue.NewNamespace("C://code/kokaq/bin", queue.NamespaceConfig{
		NamespaceName: "data-db",
		NamespaceId:   1,
	})

	logger.ConsoleLog("INFO", "Namespace created: #", ns.Id, ": ", ns.Name)

	qConfig := queue.QueueConfiguration{
		QueueName:       "test-queue",
		QueueId:         1,
		EnableDLQ:       true,
		EnableInvisible: true,
	}
	var q *queue.Queue
	var err error

	if q, err = ns.AddQueue(&qConfig); err != nil {
		logger.ConsoleLog("ERROR", "Error adding queue:", err)
	} else {
		logger.ConsoleLog("INFO", "Queue added: #", q.Id, ": ", q.Name)
	}

	// if empty, _ := q.IsEmpty(); empty {
	// 	fmt.Println("Queue is Empty")
	// }

	var qi = &queue.QueueItem{
		MessageId: uuid.New(),
		Priority:  1,
	}

	if err = q.Enqueue(qi); err != nil {
		logger.ConsoleLog("ERROR", "Error enqueuing item:", err)
	} else {
		logger.ConsoleLog("INFO", "Enqueued item:", qi.MessageId, "with priority", qi.Priority)
	}

	qi = &queue.QueueItem{
		MessageId: uuid.New(),
		Priority:  2,
	}
	if err = q.Enqueue(qi); err != nil {
		logger.ConsoleLog("ERROR", "Error enqueuing item:", err)
	} else {
		logger.ConsoleLog("INFO", "Enqueued item:", qi.MessageId, "with priority", qi.Priority)
	}

	qi = &queue.QueueItem{
		MessageId: uuid.New(),
		Priority:  2,
	}
	if err = q.Enqueue(qi); err != nil {
		logger.ConsoleLog("ERROR", "Error enqueuing item:", err)
	} else {
		logger.ConsoleLog("INFO", "Enqueued item:", qi.MessageId, "with priority", qi.Priority)
	}

	qi = &queue.QueueItem{
		MessageId: uuid.New(),
		Priority:  3,
	}
	if err = q.Enqueue(qi); err != nil {
		logger.ConsoleLog("ERROR", "Error enqueuing item:", err)
	} else {
		logger.ConsoleLog("INFO", "Enqueued item:", qi.MessageId, "with priority", qi.Priority)
	}

	if qi, err := q.Dequeue(); err != nil {
		logger.ConsoleLog("INFO", "Error dequeuing item:", err)
	} else {
		logger.ConsoleLog("INFO", "Dequeued item:", qi.MessageId, "with priority", qi.Priority)
	}

	if qi, err := q.Dequeue(); err != nil {
		logger.ConsoleLog("INFO", "Error dequeuing item:", err)
	} else {
		logger.ConsoleLog("INFO", "Dequeued item:", qi.MessageId, "with priority", qi.Priority)
	}

	if qi, err := q.Peek(); err != nil {
		logger.ConsoleLog("INFO", "Error peeking item:", err)
	} else {
		logger.ConsoleLog("INFO", "Peeked item:", qi.MessageId, "with priority", qi.Priority)
	}
	ns.DeleteQueue(1)
}

package queue

import "github.com/google/uuid"

type KokaqCore struct {
	totalNodes              int
	totalPages              int
	currentLoadedPageNumber int
	currentLoadedPage       []byte
	pagesPath               string
}

type heapNode struct {
	priority int // Priority of the node
	indexPos int // Position of the node in the heap
}

type Kokaq struct {
	namespaceId   uint32
	queueId       uint32
	pHeap         *KokaqCore
	directoryPath string
}

type QueueItem struct {
	Id       uuid.UUID
	Priority int
}

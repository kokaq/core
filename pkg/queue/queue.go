package queue

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/kokaq/core/v1/pkg/utils"
)

func NewQueueItem(id uuid.UUID, priority int) *QueueItem {
	return &QueueItem{
		Id:       id,
		Priority: priority,
	}
}

func NewKokaq(namespaceId uint32, queueId uint32) (*Kokaq, error) {

	//Check whether directory exists or not and if not create it
	dirPath := filepath.Join(STORAGE_ROOT_DIR, fmt.Sprint(namespaceId), fmt.Sprint(queueId))
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating directory:", dirPath, err)
			return nil, err
		}
	}

	//Initialize a new core heap
	heapCore, err := NewKokaqCore(filepath.Join(dirPath, STORAGE_PAGES_DIR))
	if err != nil {
		fmt.Println("Error creating core heap:", namespaceId, queueId, err)
		return nil, err
	}
	return &Kokaq{
		namespaceId:   namespaceId,
		queueId:       queueId,
		pHeap:         heapCore,
		directoryPath: dirPath,
	}, nil
}

func (pq *Kokaq) PushItem(item *QueueItem) error {
	//Priority 0 is reserved for pagging in core heap
	if item.Priority == 0 {
		errMsg := "priority 0 is not allowed"
		fmt.Println(errMsg)
		return errors.New(errMsg)
	}

	// Create index file if it doesn't exist
	indexPath := pq.getIndexFilePath(item.Priority)
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		err2 := pq.pHeap.Push(&heapNode{priority: item.Priority, indexPos: 0})
		if err2 != nil {
			fmt.Println("Error pushing item to core heap", err2)
			return err2
		}
	}
	// Put item id in index file
	byteArray := item.Id[:]
	utils.AppendBytesToFile(indexPath, byteArray)
	return nil
}

func (pq *Kokaq) PopItem() (*QueueItem, error) {

	// Peek item with highest priority
	node, err := pq.pHeap.Peek()
	if err != nil {
		fmt.Println("Error peeking item from core heap", err)
		return nil, err
	}

	priority := node.priority
	indexPath := pq.getIndexFilePath(priority)
	data, err := utils.ReadBytesFromFile(indexPath, node.indexPos*STORAGE_MESSAGE_ID_SIZE, 2*STORAGE_MESSAGE_ID_SIZE)
	if err != nil {
		fmt.Println("Error reading index file", err)
		return nil, err
	}
	itemId, err := uuid.FromBytes(data[:STORAGE_MESSAGE_ID_SIZE])
	if err != nil {
		fmt.Println("Message id is not valid in index file", err)
		return nil, err
	}
	if itemId == uuid.Nil {
		errMsg := "message id is not valid in index file"
		fmt.Println(errMsg, err)
		return nil, errors.New(errMsg)
	}

	nextItemId, err := uuid.FromBytes(data[STORAGE_MESSAGE_ID_SIZE:])
	if err != nil && nextItemId == uuid.Nil {
		// Delete index file since no emore items with that priority
		if err := os.Remove(indexPath); err != nil {
			fmt.Println("Error removing index file:", indexPath, err)
			return nil, err
		}
		_, err := pq.pHeap.Pop()
		if err != nil {
			fmt.Println("Error poping item from core heap", err)
			return nil, err
		}
	} else {
		err := pq.pHeap.SetIndexofPeekElement(node.indexPos + 1)
		if err != nil {
			fmt.Println("Error incrementing index of top element from core heap", err)
			return nil, err
		}
	}
	return NewQueueItem(itemId, priority), nil
}

func (pq *Kokaq) PeekItem() (*QueueItem, error) {

	// Peek item with highest priority
	node, err := pq.pHeap.Peek()
	if err != nil {
		fmt.Println("Error peeking item from core heap", err)
		return nil, err
	}

	priority := node.priority
	indexPath := pq.getIndexFilePath(priority)
	data, err := utils.ReadBytesFromFile(indexPath, node.indexPos*STORAGE_MESSAGE_ID_SIZE, STORAGE_MESSAGE_ID_SIZE)
	if err != nil {
		fmt.Println("Error reading index file", err)
		return nil, err
	}
	itemId, err := uuid.FromBytes(data)
	if err != nil {
		fmt.Println("Message id is not valid in index file", err)
		return nil, err
	}
	if itemId == uuid.Nil {
		errMsg := "message id is not valid in index file"
		fmt.Println(errMsg, err)
		return nil, errors.New(errMsg)
	}
	return NewQueueItem(itemId, priority), nil
}

func (pq *Kokaq) IsEmpty() bool {

	node, err := pq.pHeap.Peek()
	if err != nil {
		panic("Error peeking item from core heap")
	}

	priority := node.priority
	indexPath := pq.getIndexFilePath(priority)
	data, err := utils.ReadBytesFromFile(indexPath, node.indexPos*STORAGE_MESSAGE_ID_SIZE, STORAGE_MESSAGE_ID_SIZE)
	if err != nil {
		panic("Error reading index file")
	}
	itemId, err := uuid.FromBytes(data)
	if err != nil {
		panic("Message id is not valid in index file")
	}
	if itemId == uuid.Nil {
		return true
	}
	return false
}

func (pq *Kokaq) DeleteQueue() error {
	// Delete directory
	if err := os.RemoveAll(pq.directoryPath); err != nil {
		fmt.Println("Error removing directory:", pq.directoryPath, err)
		return err
	}
	return nil
}

func (pq *Kokaq) getIndexFilePath(priority int) string {
	return filepath.Join(pq.directoryPath, fmt.Sprintf("index-%d", priority))
}

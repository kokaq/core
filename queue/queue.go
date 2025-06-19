package queue

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/kokaq/core/utils"
)

type Kokaq struct {
	namespaceId   uint32
	queueId       uint32
	priorityHeap  *KokaqHeap
	directoryPath string
	messageIdSize int
}

func (pq *Kokaq) PushItem(item *KokaqItem) error {
	//Priority 0 is reserved for pagging in core heap
	if item.Priority == 0 {
		errMsg := "priority 0 is not allowed"
		fmt.Println(errMsg)
		return errors.New(errMsg)
	}

	// Create index file if it doesn't exist
	indexPath := pq.getIndexFilePath(item.Priority)
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		err2 := pq.priorityHeap.Push(NewHeapNode(item.Priority, 0))
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

func (pq *Kokaq) PopItem() (*KokaqItem, error) {

	// Peek item with highest priority
	node, err := pq.priorityHeap.Peek()
	if err != nil {
		fmt.Println("Error peeking item from core heap", err)
		return nil, err
	}

	priority := node.Priority
	indexPath := pq.getIndexFilePath(priority)
	data, err := utils.ReadBytesFromFile(indexPath, node.Index*pq.messageIdSize, 2*pq.messageIdSize)
	if err != nil {
		fmt.Println("Error reading index file", err)
		return nil, err
	}
	itemId, err := uuid.FromBytes(data[:pq.messageIdSize])
	if err != nil {
		fmt.Println("Message id is not valid in index file", err)
		return nil, err
	}
	if itemId == uuid.Nil {
		errMsg := "message id is not valid in index file"
		fmt.Println(errMsg, err)
		return nil, errors.New(errMsg)
	}

	nextItemId, err := uuid.FromBytes(data[pq.messageIdSize:])
	if err != nil && nextItemId == uuid.Nil {
		// Delete index file since no emore items with that priority
		if err := os.Remove(indexPath); err != nil {
			fmt.Println("Error removing index file:", indexPath, err)
			return nil, err
		}
		_, err := pq.priorityHeap.Pop()
		if err != nil {
			fmt.Println("Error poping item from core heap", err)
			return nil, err
		}
	} else {
		err := pq.priorityHeap.SetIndexofPeekElement(node.Index + 1)
		if err != nil {
			fmt.Println("Error incrementing index of top element from core heap", err)
			return nil, err
		}
	}
	return NewKokaqItem(itemId, priority), nil
}

func (pq *Kokaq) PeekItem() (*KokaqItem, error) {

	// Peek item with highest priority
	node, err := pq.priorityHeap.Peek()
	if err != nil {
		fmt.Println("Error peeking item from core heap", err)
		return nil, err
	}

	priority := node.Priority
	indexPath := pq.getIndexFilePath(priority)
	data, err := utils.ReadBytesFromFile(indexPath, node.Index*pq.messageIdSize, pq.messageIdSize)
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
	return NewKokaqItem(itemId, priority), nil
}

func (pq *Kokaq) IsEmpty() bool {

	node, err := pq.priorityHeap.Peek()
	if err != nil {
		panic("Error peeking item from core heap")
	}

	priority := node.Priority
	indexPath := pq.getIndexFilePath(priority)
	data, err := utils.ReadBytesFromFile(indexPath, node.Index*pq.messageIdSize, pq.messageIdSize)
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

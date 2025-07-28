package queue

import (
	"fmt"
	"path/filepath"

	"github.com/kokaq/core/utils"
)

type NamespaceConfig struct {
	NamespaceName string
	NamespaceId   uint32
}

type Namespace struct {
	Name    string
	Id      uint32
	RootDir string
	Queues  map[uint32]*Queue
}

func NewNamespace(parentDirectory string, config NamespaceConfig) *Namespace {
	n := &Namespace{
		Name:    config.NamespaceName,
		Queues:  make(map[uint32]*Queue, 0),
		Id:      config.NamespaceId,
		RootDir: filepath.Join(parentDirectory, fmt.Sprintf("%s-%d", config.NamespaceName, config.NamespaceId)),
	}
	if err := utils.EnsureDirectoryCreated(n.RootDir); err != nil {
		panic(fmt.Errorf("failed to create root directory for namespace %s: %w", n.Name, err))
	}
	return n
}

func (n *Namespace) GetQueue(queueId uint32) (*Queue, error) {
	if queue, exists := n.Queues[queueId]; exists {
		return queue, nil
	}
	return nil, fmt.Errorf("queue with id %d not found", queueId)
}

func (n *Namespace) AddQueue(q *QueueConfiguration) (*Queue, error) {
	var err error
	if n.Queues[q.QueueId], err = NewQueue(n.RootDir, *q); err != nil {
		return nil, fmt.Errorf("failed to add queue %s: %w", q.QueueName, err)
	}
	return n.Queues[q.QueueId], nil
}

func (n *Namespace) LoadQueue(q *QueueConfiguration) (*Queue, error) {
	if queue, exists := n.Queues[q.QueueId]; exists {
		return queue, nil
	}
	return n.AddQueue(q)
}

func (n *Namespace) ClearQueue(queueId uint32) error {
	if _, exists := n.Queues[queueId]; exists {
		if err := n.Queues[queueId].Clear(); err != nil {
			return fmt.Errorf("failed to clear queue: %w", err)
		}
	}
	return nil
}

func (n *Namespace) DeleteQueue(queueId uint32) error {
	if _, exists := n.Queues[queueId]; exists {
		n.Queues[queueId].Delete()
		n.Queues[queueId] = nil
	}
	rootDir := filepath.Join(n.RootDir, fmt.Sprint(queueId))
	if err := utils.EnsureDirectoryDeleted(rootDir); err != nil {
		return fmt.Errorf("failed to delete queue directory %s: %w", rootDir, err)
	}

	return nil
}

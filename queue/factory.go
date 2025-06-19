package queue

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/kokaq/core/utils"
)

func NewHeapNode(priority int, indexPos int) *HeapNode {
	return &HeapNode{priority, indexPos}
}

func NewKokaqHeap(config NewKokaqHeapConfiguration) (*KokaqHeap, error) {

	// SubheapLastLayerNodes = 1 << (HeapMaxSize - 1)
	// SubheapNodes          = (1 << HeapMaxSize) - 1
	// NodeSizeInBytes       = PrioritySizeInBytes + IndexSizeInBytes
	// SubheapSizeInBytes    = NodeSizeInBytes * SubheapNodes

	var subheapNodes = (1 << config.HeapMaxSize) - 1
	var nodeSizeInBytes = config.PrioritySizeInBytes + config.IndexSizeInBytes

	// if exists
	if _, err := os.Stat(config.PagesPath); os.IsNotExist(err) {
		file, err := os.Create(config.PagesPath)
		if err != nil {
			fmt.Println("Error creating pages file: ", config.PagesPath, err)
			return nil, err
		}
		defer file.Close()
		return &KokaqHeap{
			totalNodes:              0,
			totalPages:              0,
			currentLoadedPageNumber: 0,
			currentLoadedPage:       nil,
			pagesPath:               config.PagesPath,
			heapMaxSize:             config.HeapMaxSize,
			prioritySizeInBytes:     config.PrioritySizeInBytes,
			indexSizeInBytes:        config.IndexSizeInBytes,
			subheapLastLayerNodes:   1 << (config.HeapMaxSize - 1),
			subheapNodes:            subheapNodes,
			nodeSizeInBytes:         config.PrioritySizeInBytes + config.IndexSizeInBytes,
			subheapSizeInBytes:      nodeSizeInBytes * subheapNodes,
		}, err
	} else {
		//TODO: need to check whether heap is logically valid or not
		cnt := 1
		nodes := 0
		pages := 0
		for {
			currentPage, err := utils.ReadBytesFromFile(config.PagesPath, (cnt-1)*nodeSizeInBytes*subheapNodes, nodeSizeInBytes*subheapNodes)
			if err != nil {
				fmt.Println("Error reading page: ", cnt, " from file: ", config.PagesPath, err)
				return nil, err
			}
			if len(currentPage) == 0 {
				break
			}
			nodesInPage := 0
			for i := 1; i <= subheapNodes; i++ {
				currentNode, err := deserializeNode(currentPage[(i-1)*(config.PrioritySizeInBytes+config.IndexSizeInBytes):i*(config.PrioritySizeInBytes+config.IndexSizeInBytes)], config.PrioritySizeInBytes)
				if err != nil {
					fmt.Println("Error deserializing node: ", i, " from page: ", cnt, " in file: ", config.PagesPath, err)
					return nil, err
				}
				if currentNode.Priority == 0 {
					continue
				}
				if i == 1 && cnt != 1 {
					continue
				}
				nodesInPage++
			}
			nodes += nodesInPage
			if nodesInPage > 0 {
				pages++
			}
			cnt++
		}

		return &KokaqHeap{
			totalNodes:              nodes,
			totalPages:              pages,
			currentLoadedPageNumber: 0,
			currentLoadedPage:       nil,
			pagesPath:               config.PagesPath,
			heapMaxSize:             config.HeapMaxSize,
			prioritySizeInBytes:     config.PrioritySizeInBytes,
			indexSizeInBytes:        config.IndexSizeInBytes,
			subheapLastLayerNodes:   1 << (config.HeapMaxSize - 1),
			subheapNodes:            subheapNodes,
			nodeSizeInBytes:         config.PrioritySizeInBytes + config.IndexSizeInBytes,
			subheapSizeInBytes:      nodeSizeInBytes * subheapNodes,
		}, nil
	}
}

func NewKokaqItem(id uuid.UUID, priority int) *KokaqItem {
	return &KokaqItem{
		Id:       id,
		Priority: priority,
	}
}

type NewKokaqConfiguration struct {
	NamespaceId         uint32
	QueueId             uint32
	RootDir             string
	MessageIdSize       int
	HeapMaxSize         int
	PrioritySizeInBytes int
	IndexSizeInBytes    int
}

type NewKokaqHeapConfiguration struct {
	PagesPath           string
	HeapMaxSize         int
	PrioritySizeInBytes int
	IndexSizeInBytes    int
}

func NewDefaultKokaq(namespaceId uint32, queueId uint32) (*Kokaq, error) {
	// RootDir             = "./data/db"
	// MessageIdSize       = 16
	// HeapMaxSize         = 5
	// PrioritySizeInBytes = 8
	// IndexSizeInBytes    = 8
	return NewKokaq(NewKokaqConfiguration{
		NamespaceId:         namespaceId,
		QueueId:             queueId,
		RootDir:             "./data/db",
		MessageIdSize:       16,
		HeapMaxSize:         5,
		PrioritySizeInBytes: 8,
		IndexSizeInBytes:    8,
	})
}

func NewKokaq(config NewKokaqConfiguration) (*Kokaq, error) {
	var namespaceId = config.NamespaceId
	var queueId = config.QueueId
	//Check whether directory exists or not and if not create it
	dirPath := filepath.Join(config.RootDir, fmt.Sprint(namespaceId), fmt.Sprint(queueId))
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating directory:", dirPath, err)
			return nil, err
		}
	}
	//Initialize a new core heap
	heapCore, err := NewKokaqHeap(NewKokaqHeapConfiguration{
		PagesPath:           filepath.Join(dirPath, "pages"),
		HeapMaxSize:         config.HeapMaxSize,
		PrioritySizeInBytes: config.PrioritySizeInBytes,
		IndexSizeInBytes:    config.IndexSizeInBytes,
	})
	if err != nil {
		fmt.Println("Error creating core heap:", namespaceId, queueId, err)
		return nil, err
	}

	invisibilityHeap, err := NewKokaqHeap(NewKokaqHeapConfiguration{
		PagesPath:           filepath.Join(dirPath, "invisible"),
		HeapMaxSize:         config.HeapMaxSize,
		PrioritySizeInBytes: config.PrioritySizeInBytes,
		IndexSizeInBytes:    config.IndexSizeInBytes,
	})
	if err != nil {
		fmt.Println("Error creating invisibility heap:", namespaceId, queueId, err)
		return nil, err
	}
	return &Kokaq{
		namespaceId:      namespaceId,
		queueId:          queueId,
		priorityHeap:     heapCore,
		invisibilityHeap: invisibilityHeap,
		directoryPath:    dirPath,
		messageIdSize:    config.MessageIdSize,
	}, nil
}

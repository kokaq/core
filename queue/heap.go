package queue

import (
	"encoding/binary"
	"fmt"
	"math/bits"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/kokaq/core/utils"
)

type HeapConfig struct {
	PagesPath    string
	IndexPath    string
	heapMaxSize  int
	prioritySize int
	indexSize    int

	nodeSize              int
	subheapSize           int
	subheapLastLayerNodes int
	subheapNodes          int
	messageIdSize         int
}

type Heap struct {
	totalNodes  int
	totalPages  int
	config      HeapConfig
	currentPage *Page
}

func NewHeap(parentDirectory string, heapMaxSize int, prioritySize int, indexSize int, messageIdSize int) (*Heap, error) {
	var err error
	var subheapNodes = (1 << heapMaxSize) - 1
	indexpath := filepath.Join(parentDirectory, "indexes") // directory
	pagesPath := filepath.Join(parentDirectory, "pages")   // file
	if err = utils.EnsureDirectoryCreated(indexpath); err != nil {
		return nil, fmt.Errorf("failed to create index directory %s: %w", indexpath, err)
	}
	if !utils.FileExists(pagesPath) {
		if err = utils.EnsureFileCreated(pagesPath); err != nil {
			return nil, fmt.Errorf("failed to create file %s: %w", pagesPath, err)
		}
		return &Heap{
			totalNodes:  0,
			totalPages:  0,
			currentPage: NewPage(),
			config: HeapConfig{
				PagesPath:             pagesPath,
				IndexPath:             indexpath,
				messageIdSize:         messageIdSize,
				heapMaxSize:           heapMaxSize,
				prioritySize:          prioritySize,
				indexSize:             indexSize,
				nodeSize:              prioritySize + indexSize,
				subheapSize:           (prioritySize + indexSize) * subheapNodes,
				subheapLastLayerNodes: (1 << (heapMaxSize - 1)),
				subheapNodes:          subheapNodes,
			},
		}, nil
	} else {
		cnt := 1
		nodes := 0
		pages := 0

		for {
			currPage, err := utils.ReadBytesFromFile(pagesPath, int64((cnt-1)*(prioritySize+indexSize)*subheapNodes), (prioritySize+indexSize)*subheapNodes)
			if err != nil {
				return nil, fmt.Errorf("failed to read file %s: %w", pagesPath, err)
			}
			if len(currPage) == 0 {
				break
			}
			nodesInPage := 0
			for i := 1; i <= subheapNodes; i++ {
				startIndex := (i - 1) * (prioritySize + indexSize)
				length := (prioritySize + indexSize)
				priority := int(binary.LittleEndian.Uint64(currPage[startIndex:length][:prioritySize]))
				if priority == 0 || (i == 1 && cnt != 1) {
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

		return &Heap{
			totalNodes:  nodes,
			totalPages:  pages,
			currentPage: NewPage(),
			config: HeapConfig{
				PagesPath:             pagesPath,
				IndexPath:             indexpath,
				heapMaxSize:           heapMaxSize,
				prioritySize:          prioritySize,
				indexSize:             indexSize,
				nodeSize:              prioritySize + indexSize,
				subheapSize:           (prioritySize + indexSize) * subheapNodes,
				subheapLastLayerNodes: (1 << (heapMaxSize - 1)),
				subheapNodes:          subheapNodes,
			},
		}, nil
	}
}

// Public Methods

func (h *Heap) Enqueue(queueItem *QueueItem) error {
	if queueItem.Priority == 0 {
		return fmt.Errorf("priority cannot be zero")
	}

	indexPath := h.getIndexFilePath(queueItem.Priority)
	if !utils.FileExists(indexPath) {
		// if this is a new priority
		// - add priority node to heap
		// - heapify up
		if h.totalNodes == 0 {
			h.loadPage(0)
			newPage := make([]byte, h.config.subheapSize)
			binary.LittleEndian.PutUint64(newPage[0:], uint64(queueItem.Priority))
			binary.LittleEndian.PutUint64(newPage[h.config.prioritySize:], uint64(0))
			h.currentPage.SetData(1, newPage)
			h.totalPages += 1
		} else {
			if h.totalNodes >= h.config.subheapLastLayerNodes {
				return fmt.Errorf("heap is full, cannot enqueue more nodes")
			}
			h.heapifyUp(h.totalNodes+1, queueItem.Priority, 0)
		}
		h.totalNodes++
	}
	byteArray := queueItem.MessageId[:]
	utils.AppendBytesToFile(indexPath, byteArray)
	return nil
}

func (h *Heap) Dequeue() (*QueueItem, error) {
	if err := h.loadPage(1); err != nil {
		return nil, fmt.Errorf("failed to load page 1: %w", err)
	}
	if len(h.currentPage.data) < (h.config.prioritySize + h.config.indexSize) {
		return nil, fmt.Errorf("heap is empty")
	}
	priority := binary.LittleEndian.Uint64(h.currentPage.data[:h.config.prioritySize])
	indexPos := binary.LittleEndian.Uint64(h.currentPage.data[h.config.prioritySize:])
	indexPath := h.getIndexFilePath(priority)
	offset := int64(indexPos) * int64(h.config.messageIdSize)
	data, err := utils.ReadBytesFromFile(indexPath, offset, 2*h.config.messageIdSize)
	if err != nil {
		return nil, fmt.Errorf("failed to read message id from index file: %w", err)
	}
	if itemId, err := uuid.FromBytes(data[:h.config.messageIdSize]); err != nil {
		return nil, fmt.Errorf("failed to convert bytes to UUID: %w", err)
	} else {
		if itemId == uuid.Nil {
			return nil, fmt.Errorf("no item found for the given priority")
		}
		nextItemId, err := uuid.FromBytes(data[h.config.messageIdSize:])
		if err != nil && nextItemId == uuid.Nil {
			if err := utils.EnsureFileDeleted(indexPath); err != nil {
				return nil, fmt.Errorf("failed to delete index file %s: %w", indexPath, err)
			}
			// core
			if h.totalNodes == 0 {
				return nil, fmt.Errorf("no items in the heap")
			}
			if h.totalNodes == 1 {
				if err = h.loadPage(1); err != nil {
					return nil, fmt.Errorf("failed to load page 1: %w", err)
				}
				// priority1 := binary.LittleEndian.Uint64(h.currentPage.data[:h.config.prioritySize])
				// indexPos1 := binary.LittleEndian.Uint64(h.currentPage.data[h.config.prioritySize:])
				copy(h.currentPage.data[:h.config.nodeSize], make([]byte, h.config.nodeSize))
				h.totalNodes -= 1
				h.totalPages -= 1
			} else {
				if err := h.loadPage(1); err != nil {
					return nil, fmt.Errorf("failed to load page 1: %w", err)
				}
				priority = binary.LittleEndian.Uint64(h.currentPage.data[:h.config.prioritySize])
				indexPos = binary.LittleEndian.Uint64(h.currentPage.data[h.config.prioritySize:])

				pageNumber, localIndex, _ := h.getLocalHeapDetailsForNode(h.totalNodes)
				if err = h.loadPage(pageNumber); err != nil {
					return nil, fmt.Errorf("failed to load page %d: %w", pageNumber, err)
				}

				startIndex := (localIndex - 1) * h.config.nodeSize
				lpriority := binary.LittleEndian.Uint64(h.currentPage.data[startIndex : startIndex+h.config.prioritySize])
				lindexPos := binary.LittleEndian.Uint64(h.currentPage.data[startIndex+h.config.prioritySize : startIndex+h.config.nodeSize])
				copy(h.currentPage.data[startIndex:], make([]byte, h.config.nodeSize))
				h.totalNodes -= 1
				if localIndex == 2 && pageNumber != 1 {
					h.totalPages -= 1
					copy(h.currentPage.data[:h.config.nodeSize], make([]byte, h.config.nodeSize))
				}
				if err = h.heapifyDown(lpriority, lindexPos); err != nil {
					return nil, fmt.Errorf("failed to heapify down after dequeue: %w", err)
				}

				return &QueueItem{
					MessageId: itemId,
					Priority:  priority,
				}, nil
			}
		} else {
			if err = h.setIndexOfPeekElement(int(indexPos + 1)); err != nil {
				return nil, fmt.Errorf("failed to set index of peek element: %w", err)
			}
		}

		return &QueueItem{
			MessageId: itemId,
			Priority:  priority,
		}, nil
	}
}

func (h *Heap) IsEmpty() (bool, error) {
	if err := h.loadPage(1); err != nil {
		return true, nil
	}
	var priority uint64 = 0
	var indexPos uint64 = 0

	indexPath := h.getIndexFilePath(priority)
	if len(h.currentPage.data) == 0 || !utils.FileExists(indexPath) {
		return true, nil
	}
	priority = binary.LittleEndian.Uint64(h.currentPage.data[:h.config.prioritySize])
	indexPos = binary.LittleEndian.Uint64(h.currentPage.data[h.config.prioritySize:])
	var offset = int64(indexPos) * int64(h.config.messageIdSize)
	data, err := utils.ReadBytesFromFile(indexPath, offset, h.config.messageIdSize)
	if err != nil {
		return false, fmt.Errorf("error reading indexfile")
	}
	itemId, err := uuid.FromBytes(data)
	if err != nil {
		return false, fmt.Errorf("messageid is not valid")
	}
	if itemId == uuid.Nil {
		return true, nil
	}
	return false, nil
}

func (h *Heap) Peek() (*QueueItem, error) {
	if err := h.loadPage(1); err != nil {
		return nil, fmt.Errorf("failed to load page 1: %w", err)
	}
	if len(h.currentPage.data) < (h.config.prioritySize + h.config.indexSize) {
		return nil, fmt.Errorf("heap is empty")
	}
	var priority uint64 = 0
	var indexPos uint64 = 0
	priority = binary.LittleEndian.Uint64(h.currentPage.data[:h.config.prioritySize])
	indexPos = binary.LittleEndian.Uint64(h.currentPage.data[h.config.prioritySize:])
	indexPath := h.getIndexFilePath(priority)
	data, err := utils.ReadBytesFromFile(indexPath, int64(indexPos)*int64(h.config.messageIdSize), h.config.messageIdSize)
	if err != nil {
		return nil, fmt.Errorf("failed to read message id from index file: %w", err)
	}
	if itemId, err := uuid.FromBytes(data); err != nil {
		return nil, fmt.Errorf("failed to convert bytes to UUID: %w", err)
	} else {
		if itemId == uuid.Nil {
			return nil, fmt.Errorf("no item found for the given priority")
		}
		queueItem := &QueueItem{
			MessageId: itemId,
			Priority:  priority,
		}
		return queueItem, nil
	}
}

func (h *Heap) GetConfig() HeapConfig {
	return h.config
}

// Internal Methods

func (h *Heap) setIndexOfPeekElement(index int) error {
	if err := h.loadPage(1); err != nil {
		return fmt.Errorf("failed to load page 1: %w", err)
	}
	binary.LittleEndian.PutUint64(h.currentPage.data[h.config.prioritySize:], uint64(index))
	return nil
}

func (h *Heap) getIndexFilePath(priority uint64) string {
	return filepath.Join(h.config.IndexPath, fmt.Sprint(priority))
}

func (h *Heap) loadPage(pageNumber int) error {
	if pageNumber < 0 {
		return fmt.Errorf("invalid page number: %d", pageNumber)
	}
	if pageNumber != h.currentPage.index {
		h.commitCurrentPage()
		return h.loadIntoCurrentPage(pageNumber)
	}
	return nil
}

func (h *Heap) loadIntoCurrentPage(pageNumber int) error {
	if pageNumber != 0 {
		startIndex := (pageNumber - 1) * (h.config.subheapSize)
		var err error
		if h.currentPage.data, err = utils.ReadBytesFromFile(h.config.PagesPath, int64(startIndex), h.config.subheapSize); err != nil {
			return fmt.Errorf("failed to read page %d: %w", pageNumber, err)
		}
		h.currentPage.index = pageNumber
	} else {
		h.currentPage.data = nil
		h.currentPage.index = 0
	}
	return nil
}

func (h *Heap) commitCurrentPage() error {
	if h.currentPage.index != 0 {
		startIndex := (h.currentPage.index - 1) * (h.config.subheapSize)
		if err := utils.WriteBytesToFile(h.config.PagesPath, int64(startIndex), h.currentPage.data); err != nil {
			return fmt.Errorf("failed to commit page %d: %w", h.currentPage.index, err)
		}
	}
	return nil
}

func (h *Heap) heapifyUp(heapIndex int, priority uint64, index uint64) error {
	pageNumber, localIndex, _ := h.getLocalHeapDetailsForNode(heapIndex)
	if localIndex == 2 && pageNumber != 1 {
		// This is the first element in the new page,
		// root node will be there in the parent page,
		// load it and then move forward
		parentIndex := heapIndex / 2
		pageOfParent, localIndexOfParent, _ := h.getLocalHeapDetailsForNode(parentIndex)
		if err := h.loadPage(pageOfParent); err != nil {
			return fmt.Errorf("failed to load parent page %d: %w", pageOfParent, err)
		}
		rootNode := h.currentPage.data[(localIndexOfParent-1)*(h.config.nodeSize) : (localIndexOfParent)*(h.config.nodeSize)]
		if err := h.loadPage(0); err != nil {
			return fmt.Errorf("failed to load current page %d: %w", 0, err)
		}
		newPage := make([]byte, h.config.subheapSize)
		copy(newPage[:h.config.nodeSize], rootNode)
		h.currentPage.SetData(pageNumber, newPage)
	}

	//Fake conditions to perform first executions irrespective of anything
	needToGoUp := true
	previousPage := 0
	for needToGoUp {
		pageNumber, localIndex, localLevel := h.getLocalHeapDetailsForNode(heapIndex)
		if err := h.loadPage(pageNumber); err != nil {
			return fmt.Errorf("failed to load page %d: %w", pageNumber, err)
		}
		startIndex := (localIndex - 1) * h.config.nodeSize

		newPage := make([]byte, h.config.nodeSize)
		binary.LittleEndian.PutUint64(newPage[0:], uint64(priority))
		binary.LittleEndian.PutUint64(newPage[h.config.prioritySize:], uint64(index))
		copy(h.currentPage.data[startIndex:], newPage)

		needToGoUp = h.localHeapifyUp(localIndex)

		if previousPage != 0 {
			offset := int64((previousPage - 1) * (h.config.subheapSize))
			_data := h.currentPage.data[startIndex : startIndex+h.config.nodeSize]
			utils.WriteBytesToFile(h.config.PagesPath, offset, _data)
		}
		if pageNumber == 1 {
			// If first sub heap no need to go up
			break
		}
		//Change parametes for next iterations
		heapIndex = heapIndex >> localLevel
		previousPage = pageNumber
	}
	return nil
}

func (h *Heap) localHeapifyUp(index int) bool {
	indexChild := index
	for indexChild > 1 {
		startPointChild := (indexChild - 1) * (h.config.nodeSize)
		bytesChild := h.currentPage.data[startPointChild : startPointChild+h.config.nodeSize]
		priorityChild := binary.LittleEndian.Uint64(bytesChild[:h.config.prioritySize])

		indexParent := indexChild / 2

		startPointParent := (indexParent - 1) * (h.config.nodeSize)
		bytesParent := h.currentPage.data[startPointParent : startPointParent+h.config.nodeSize]
		priorityParent := binary.LittleEndian.Uint64(bytesParent[:h.config.prioritySize])

		// TODO: Priority comparison should be configurable
		if priorityChild > priorityParent {
			tempByteParents := make([]byte, len(bytesParent))
			copy(tempByteParents, bytesParent)
			copy(h.currentPage.data[startPointParent:], bytesChild)
			copy(h.currentPage.data[startPointChild:], tempByteParents)
		} else {
			return false
		}
		indexChild = indexChild / 2
	}
	return true
}

func (h *Heap) heapifyDown(priority uint64, index uint64) error {
	needToGoDown := true
	pageNumber := 1
	previousPage := 0
	previousIndex := 0
	lastLayerIndex := 0
	for needToGoDown {
		var err error
		if err = h.loadPage(pageNumber); err != nil {
			return fmt.Errorf("failed to load page %d: %w", pageNumber, err)
		}

		binary.LittleEndian.PutUint64(h.currentPage.data[:h.config.prioritySize], priority)
		binary.LittleEndian.PutUint64(h.currentPage.data[h.config.prioritySize:], index)
		if needToGoDown, lastLayerIndex, err = h.localHeapifyDown(); err != nil {
			return fmt.Errorf("failed to heapify down: %w", err)
		}

		//If second iteration change the last layer node of previous page
		//------Anti-pattern code------
		// Here we are directly going to previous page and changing the last layer node,
		// every operation should be done by loading page
		// Here, we are sure that page is synced in meomry,
		//  to reduce disk I/O this has been done

		if previousPage != 0 {
			var offset = int64((previousPage-1)*(h.config.subheapSize) + (previousIndex-1)*(h.config.nodeSize))
			utils.WriteBytesToFile(h.config.PagesPath, offset, h.currentPage.data[:h.config.nodeSize])
		}
		previousPage = pageNumber
		previousIndex = lastLayerIndex

		pageNumber = (h.config.subheapLastLayerNodes)*(pageNumber-1) + (lastLayerIndex - h.config.subheapLastLayerNodes + 1) + 1

		if pageNumber > h.totalPages {
			break
		}
	}
	return nil
}

func (h *Heap) localHeapifyDown() (bool, int, error) {
	indexParent := 1

	for indexParent < h.config.subheapLastLayerNodes {
		parentBytes := h.currentPage.data[(indexParent-1)*(h.config.nodeSize) : (indexParent)*(h.config.nodeSize)]
		priorityParent := binary.LittleEndian.Uint64(parentBytes[:h.config.prioritySize])
		indexLeftChild := indexParent * 2
		indexRightChild := indexLeftChild + 1
		leftChildBytes := h.currentPage.data[(indexLeftChild-1)*(h.config.nodeSize) : (indexLeftChild)*(h.config.nodeSize)]
		rightChildBytes := h.currentPage.data[(indexRightChild-1)*(h.config.nodeSize) : (indexRightChild)*(h.config.nodeSize)]
		priorityLeftChild := binary.LittleEndian.Uint64(leftChildBytes[:h.config.prioritySize])
		priorityRightChild := binary.LittleEndian.Uint64(rightChildBytes[:h.config.prioritySize])

		//None of the childern exist, so no need to go further
		if priorityLeftChild == 0 && priorityRightChild == 0 {
			return false, indexParent, nil
		}

		//Right child does not exist
		if priorityRightChild == 0 {
			//Parent wins, no need to go further
			// TODO: Priority comparison should be configurable
			if priorityParent > priorityLeftChild {
				return false, indexParent, nil
				//Left child wins still no need to go further
			} else {
				tempByteParents := make([]byte, len(parentBytes))
				copy(tempByteParents, parentBytes)
				copy(h.currentPage.data[(indexParent-1)*(h.config.nodeSize):], leftChildBytes)
				copy(h.currentPage.data[(indexLeftChild-1)*(h.config.nodeSize):], tempByteParents)
				return false, indexLeftChild, nil
			}
		} else {
			//Both children exist
			// TODO: Priority comparison should be configurable
			if priorityParent > priorityLeftChild && priorityParent > priorityRightChild {
				//parent wins, no need to go further
				return false, indexParent, nil
			} else {
				//Left child wins
				// TODO: Priority comparison should be configurable
				if priorityLeftChild > priorityRightChild {
					tempByteParents := make([]byte, len(parentBytes))
					copy(tempByteParents, parentBytes)
					copy(h.currentPage.data[(indexParent-1)*(h.config.nodeSize):], leftChildBytes)
					copy(h.currentPage.data[(indexLeftChild-1)*(h.config.nodeSize):], tempByteParents)
					indexParent = indexLeftChild
				} else {
					//Right child wins
					tempByteParents := make([]byte, len(parentBytes))
					copy(tempByteParents, parentBytes)
					copy(h.currentPage.data[(indexParent-1)*(h.config.nodeSize):], rightChildBytes)
					copy(h.currentPage.data[(indexRightChild-1)*(h.config.nodeSize):], tempByteParents)
					indexParent = indexRightChild
				}
			}
		}
	}
	return true, indexParent, nil
}

func (h *Heap) getLocalHeapDetailsForNode(index int) (int, int, int) {
	globalLevel := bits.Len(uint(index)) - 1
	nodesInGloablLevel := 1 << globalLevel
	subHeapLevel := (globalLevel - 1) / (h.config.heapMaxSize - 1)
	localLevel := ((globalLevel - 1) % (h.config.heapMaxSize - 1)) + 1
	nodesInLocalLevel := 1 << localLevel
	pPartial := (index - nodesInGloablLevel) / (nodesInLocalLevel)
	pFull := (utils.Power(h.config.subheapLastLayerNodes, subHeapLevel) - 1) / (h.config.subheapLastLayerNodes - 1)
	p := pFull + pPartial + 1
	localIndexFull := nodesInLocalLevel - 1
	localIndexPartial := (index - nodesInGloablLevel) % nodesInLocalLevel
	localIndex := localIndexFull + localIndexPartial + 1
	return p, localIndex, localLevel
}

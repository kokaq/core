package queue

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math/bits"

	"github.com/kokaq/core/utils"
)

type KokaqHeap struct {
	totalNodes              int
	totalPages              int
	currentLoadedPageNumber int
	currentLoadedPage       []byte
	pagesPath               string
	heapMaxSize             int
	prioritySizeInBytes     int
	indexSizeInBytes        int
	subheapLastLayerNodes   int
	subheapNodes            int
	nodeSizeInBytes         int
	subheapSizeInBytes      int
}

// Heap Public Functions

func (queue *KokaqHeap) Push(node *HeapNode) error {
	if queue.totalNodes == 0 {
		queue.loadPage(0)
		newPage := make([]byte, queue.subheapSizeInBytes)
		newNode, err := serializeNode(node, queue.prioritySizeInBytes, queue.indexSizeInBytes, queue.nodeSizeInBytes)
		if err != nil {
			fmt.Println("Error serializing node", err)
			return err
		}
		copy(newPage[:queue.nodeSizeInBytes], newNode)
		queue.saveNewPage(1, newPage)
	} else {
		queue.pageHeapifyUp(queue.totalNodes+1, node)
	}
	queue.totalNodes += 1
	return nil
}

func (queue *KokaqHeap) Pop() (*HeapNode, error) {
	if queue.totalNodes == 0 {
		errMsg := "no nodes to pop"
		fmt.Println(errMsg)
		return nil, errors.New(errMsg)
	}
	if queue.totalNodes == 1 {
		page, err := queue.loadPage(1)
		if err != nil {
			return nil, err
		}

		newEle := make([]byte, queue.nodeSizeInBytes)
		oldEle, err := deserializeNode(page[:queue.nodeSizeInBytes], queue.prioritySizeInBytes)
		if err != nil {
			return nil, err
		}
		copy(page[:queue.nodeSizeInBytes], newEle)
		queue.totalNodes -= 1
		queue.totalPages -= 1
		return oldEle, nil
	} else {
		//Get top node
		topNode, err := queue.Peek()
		if err != nil {
			return nil, err
		}
		//Get last node and replace it with all zeros
		pageNumber, localIndex, _ := queue.getLocalHeapDetailsForNode(queue.totalNodes)
		page, err := queue.loadPage(pageNumber)
		if err != nil {
			return nil, err
		}
		startPoint := (localIndex - 1) * (queue.nodeSizeInBytes)
		lastNode, err := deserializeNode(page[startPoint:startPoint+queue.nodeSizeInBytes], queue.prioritySizeInBytes)
		if err != nil {
			return nil, err
		}
		newNode := make([]byte, queue.nodeSizeInBytes)
		copy(page[startPoint:], newNode)
		queue.totalNodes -= 1
		if localIndex == 2 && pageNumber != 1 {
			queue.totalPages -= 1
			//Remove the root node as well
			copy(page[:queue.nodeSizeInBytes], newNode)
		}
		err = queue.pageHeapifyDown(lastNode)
		if err != nil {
			return nil, err
		}
		return topNode, nil
	}
}

func (queue *KokaqHeap) Peek() (*HeapNode, error) {
	page, err := queue.loadPage(1)
	if err != nil {
		return nil, err
	}
	return deserializeNode(page[:queue.nodeSizeInBytes], queue.prioritySizeInBytes)
}

func (queue *KokaqHeap) SetIndexofPeekElement(index int) error {
	page, err := queue.loadPage(1)
	if err != nil {
		return err
	}
	node, err := deserializeNode(page[:queue.nodeSizeInBytes], queue.prioritySizeInBytes)
	if err != nil {
		return err
	}

	node.Index = index
	newNode, err := serializeNode(node, queue.prioritySizeInBytes, queue.indexSizeInBytes, queue.nodeSizeInBytes)
	if err != nil {
		return err
	}

	copy(page[:queue.nodeSizeInBytes], newNode)
	return nil
}

// Heap Private Functions

func (queue *KokaqHeap) pageHeapifyUp(index int, node *HeapNode) error {
	pageNumber, localIndex, _ := queue.getLocalHeapDetailsForNode(index)
	if localIndex == 2 && pageNumber != 1 {
		//This is the first element in the new page, root node will be there in the parent page, load it and then move forward
		parentIndex := index / 2
		pageOfParent, localIndexOfParent, _ := queue.getLocalHeapDetailsForNode(parentIndex)
		page, err := queue.loadPage(pageOfParent)
		if err != nil {
			fmt.Println("Error loading page: ", pageOfParent, " in file: ", queue.pagesPath, err)
			return err
		}
		rootNode := page[(localIndexOfParent-1)*(queue.nodeSizeInBytes) : (localIndexOfParent)*(queue.nodeSizeInBytes)]
		_, err = queue.loadPage(0)
		if err != nil {
			fmt.Println("Error loading page: 0 in file: ", queue.pagesPath, err)
			return err
		}
		newPage := make([]byte, queue.subheapSizeInBytes)
		copy(newPage[:queue.nodeSizeInBytes], rootNode)
		queue.saveNewPage(pageNumber, newPage)
	}

	//Fake conditions to perform first executions irrespective of anything
	needToGoUp := true
	previousPage := 0
	for needToGoUp {
		pageNumber, localIndex, localLevel := queue.getLocalHeapDetailsForNode(index)
		page, err := queue.loadPage(pageNumber)
		if err != nil {
			fmt.Println("Error loading page: ", pageNumber, " in file: ", queue.pagesPath, err)
			return err
		}
		startPoint := (localIndex - 1) * (queue.nodeSizeInBytes)
		nodeContent, err := serializeNode(node, queue.prioritySizeInBytes, queue.indexSizeInBytes, queue.nodeSizeInBytes)
		if err != nil {
			fmt.Println("Error serializing node", err)
			return err
		}
		copy(page[startPoint:], nodeContent)
		needToGoUp, err = queue.localHeapifyUp(localIndex, page)
		if err != nil {
			return err
		}

		//Change the root node of previous page
		//------Anti-pattern code------
		// Here we are directly going to previous page and changing the root node, every operation should be doen by loading page
		// Here, we are sure that page is synced in meomry, to reduce disk I/O this has been done
		if previousPage != 0 {
			utils.WriteBytesToFile(queue.pagesPath, int64((previousPage-1)*(queue.subheapSizeInBytes)), page[startPoint:startPoint+queue.nodeSizeInBytes])
		}
		//If first sub heap no need to go up
		if pageNumber == 1 {
			break
		}
		//Change parametes for next iterations
		index = index >> localLevel
		previousPage = pageNumber
	}
	return nil
}

func (queue *KokaqHeap) localHeapifyUp(index int, heap []byte) (bool, error) {
	indexChild := index
	for indexChild > 1 {
		startPointChild := (indexChild - 1) * (queue.nodeSizeInBytes)
		bytesChild := heap[startPointChild : startPointChild+queue.nodeSizeInBytes]
		priorityChild := binary.LittleEndian.Uint64(bytesChild[:queue.prioritySizeInBytes])

		indexParent := indexChild / 2
		startPointParent := (indexParent - 1) * (queue.nodeSizeInBytes)
		bytesParent := heap[startPointParent : startPointParent+queue.nodeSizeInBytes]
		priorityParent := binary.LittleEndian.Uint64(bytesParent[:queue.prioritySizeInBytes])

		if priorityChild > priorityParent {
			tempByteParents := make([]byte, len(bytesParent))
			copy(tempByteParents, bytesParent)
			copy(heap[startPointParent:], bytesChild)
			copy(heap[startPointChild:], tempByteParents)
		} else {
			return false, nil
		}
		indexChild = indexChild / 2
	}
	return true, nil
}

func (queue *KokaqHeap) pageHeapifyDown(node *HeapNode) error {
	needToGoDown := true
	pageNumber := 1
	previousPage := 0
	previousIndex := 0
	lastLayerIndex := 0
	for needToGoDown {
		page, err := queue.loadPage(pageNumber)
		if err != nil {
			return err
		}
		newNode, err := serializeNode(node, queue.prioritySizeInBytes, queue.indexSizeInBytes, queue.nodeSizeInBytes)
		if err != nil {
			return err
		}
		copy(page[:queue.nodeSizeInBytes], newNode)
		needToGoDown, lastLayerIndex, err = queue.localHeapifyDown(page)
		if err != nil {
			return err
		}
		//If second iteration change the last layer node of previous page
		//------Anti-pattern code------
		// Here we are directly going to previous page and changing the last layer node, every operation should be doen by loading page
		// Here, we are sure that page is synced in meomry, to reduce disk I/O this has been done
		if previousPage != 0 {
			utils.WriteBytesToFile(queue.pagesPath, int64((previousPage-1)*(queue.subheapSizeInBytes)+(previousIndex-1)*(queue.nodeSizeInBytes)), page[:queue.nodeSizeInBytes])
		}

		//Change parametes for next iterations
		previousPage = pageNumber
		previousIndex = lastLayerIndex

		//Next page number which needs to be loaded based on current page we are on and last layer index
		pageNumber = (queue.subheapLastLayerNodes)*(pageNumber-1) + (lastLayerIndex - queue.subheapLastLayerNodes + 1) + 1

		//If no pages available, no need to go down
		if pageNumber > queue.totalPages {
			break
		}
	}
	return nil
}

func (queue *KokaqHeap) localHeapifyDown(heap []byte) (bool, int, error) {
	indexParent := 1
	for indexParent < queue.subheapLastLayerNodes {
		parentBytes := heap[(indexParent-1)*(queue.nodeSizeInBytes) : (indexParent)*(queue.nodeSizeInBytes)]
		priorityParent := binary.LittleEndian.Uint64(parentBytes[:queue.prioritySizeInBytes])
		indexLeftChild := indexParent * 2
		indexRightChild := indexLeftChild + 1
		leftChildBytes := heap[(indexLeftChild-1)*(queue.nodeSizeInBytes) : (indexLeftChild)*(queue.nodeSizeInBytes)]
		rightChildBytes := heap[(indexRightChild-1)*(queue.nodeSizeInBytes) : (indexRightChild)*(queue.nodeSizeInBytes)]
		priorityLeftChild := binary.LittleEndian.Uint64(leftChildBytes[:queue.prioritySizeInBytes])
		priorityRightChild := binary.LittleEndian.Uint64(rightChildBytes[:queue.prioritySizeInBytes])

		//None of the childern exist, so no need to go further
		if priorityLeftChild == 0 && priorityRightChild == 0 {
			return false, indexParent, nil
		}

		//Right child does not exist
		if priorityRightChild == 0 {
			//Parent wins, no need to go further
			if priorityParent > priorityLeftChild {
				return false, indexParent, nil
				//Left child wins still no need to go further
			} else {
				tempByteParents := make([]byte, len(parentBytes))
				copy(tempByteParents, parentBytes)
				copy(heap[(indexParent-1)*(queue.nodeSizeInBytes):], leftChildBytes)
				copy(heap[(indexLeftChild-1)*(queue.nodeSizeInBytes):], tempByteParents)
				return false, indexLeftChild, nil
			}
		} else {
			//Both children exist
			if priorityParent > priorityLeftChild && priorityParent > priorityRightChild {
				//parent wins, no need to go further
				return false, indexParent, nil
			} else {
				//Left child wins
				if priorityLeftChild > priorityRightChild {
					tempByteParents := make([]byte, len(parentBytes))
					copy(tempByteParents, parentBytes)
					copy(heap[(indexParent-1)*(queue.nodeSizeInBytes):], leftChildBytes)
					copy(heap[(indexLeftChild-1)*(queue.nodeSizeInBytes):], tempByteParents)
					indexParent = indexLeftChild
				} else {
					//Right child wins
					tempByteParents := make([]byte, len(parentBytes))
					copy(tempByteParents, parentBytes)
					copy(heap[(indexParent-1)*(queue.nodeSizeInBytes):], rightChildBytes)
					copy(heap[(indexRightChild-1)*(queue.nodeSizeInBytes):], tempByteParents)
					indexParent = indexRightChild
				}
			}
		}
	}
	return true, indexParent, nil
}

// Disk IO functions
// TODO: Mark it as critical section
// Load page will first check whether that page is there is the cached property or not
// If yes it will not make disk IO, else it will commit the previous page and load requested page from disk
// If pageNumber is 0, it's special case: it will not load any page and will just commit the current page.
// It's used to free up the space so before creating any new page in memory loadPage(0) should be called
func (queue *KokaqHeap) loadPage(pageNumber int) ([]byte, error) {
	if pageNumber < 0 {
		errMsg := "page number cannot be negative"
		fmt.Println(errMsg)
		return nil, errors.New(errMsg)
	}

	//page is already loaded
	if queue.currentLoadedPageNumber == pageNumber {
		return queue.currentLoadedPage, nil
	}

	//page is not loaded
	//commit the page if it's not zero, since we are loading a new page
	if queue.currentLoadedPageNumber != 0 {
		startIndex := (queue.currentLoadedPageNumber - 1) * (queue.subheapSizeInBytes)
		err := utils.WriteBytesToFile(queue.pagesPath, int64(startIndex), queue.currentLoadedPage)
		if err != nil {
			fmt.Println("Error saving page: ", queue.currentLoadedPageNumber, " to file: ", queue.pagesPath, err)
			return nil, err
		}
	}
	if pageNumber != 0 {
		startIndex := (pageNumber - 1) * (queue.subheapSizeInBytes)
		var err error
		queue.currentLoadedPage, err = utils.ReadBytesFromFile(queue.pagesPath, startIndex, queue.subheapSizeInBytes)
		if err != nil {
			fmt.Println("Error loading page: ", pageNumber, " from file: ", queue.pagesPath, err)
			return nil, err
		}
		queue.currentLoadedPageNumber = pageNumber
	} else {
		queue.currentLoadedPage = nil
		queue.currentLoadedPageNumber = 0
	}
	return queue.currentLoadedPage, nil
}

// It will save the given page (THIS SAVE IS IN MEMORY ONLY, NOT ON DISK)
// TODO: Mark it as critical section
// It's caller's responsibility to first save the original page, then create a new page and commit
// loadPage(0)          --> required to free up the space in RAM
// create new page
// commit new page
func (queue *KokaqHeap) saveNewPage(pageNumber int, data []byte) error {
	if pageNumber < 1 {
		errMsg := "page number cannot be less than 1"
		fmt.Println(errMsg)
		return errors.New(errMsg)
	}
	queue.currentLoadedPage = data
	queue.currentLoadedPageNumber = pageNumber
	queue.totalPages += 1
	return nil
}

// Mapping functions
func (queue *KokaqHeap) getLocalHeapDetailsForNode(index int) (int, int, int) {
	globalLevel := bits.Len(uint(index)) - 1
	nodesInGloablLevel := 1 << globalLevel
	subHeapLevel := (globalLevel - 1) / (queue.heapMaxSize - 1)
	localLevel := ((globalLevel - 1) % (queue.heapMaxSize - 1)) + 1
	nodesInLocalLevel := 1 << localLevel
	pPartial := (index - nodesInGloablLevel) / (nodesInLocalLevel)
	pFull := (utils.Power(queue.subheapLastLayerNodes, subHeapLevel) - 1) / (queue.subheapLastLayerNodes - 1)
	p := pFull + pPartial + 1
	localIndexFull := nodesInLocalLevel - 1
	localIndexPartial := (index - nodesInGloablLevel) % nodesInLocalLevel
	localIndex := localIndexFull + localIndexPartial + 1
	return p, localIndex, localLevel
}

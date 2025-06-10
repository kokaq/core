package queue

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math/bits"
	"os"

	"github.com/kokaq/core/v1/pkg/utils"
)

func NewKokaqCore(pagesPath string) (*KokaqCore, error) {
	// if exists
	if _, err := os.Stat(pagesPath); os.IsNotExist(err) {
		file, err := os.Create(pagesPath)
		if err != nil {
			fmt.Println("Error creating pages file: ", pagesPath, err)
			return nil, err
		}
		defer file.Close()
		return &KokaqCore{
			totalNodes:              0,
			totalPages:              0,
			currentLoadedPageNumber: 0,
			currentLoadedPage:       nil,
			pagesPath:               pagesPath,
		}, err
	} else {
		//TODO: need to check whether heap is logically valid or not
		cnt := 1
		nodes := 0
		pages := 0
		for {
			currentPage, err := utils.ReadBytesFromFile(pagesPath, (cnt-1)*STORAGE_SUBHEAP_SIZE_IN_BYTES, STORAGE_SUBHEAP_SIZE_IN_BYTES)
			if err != nil {
				fmt.Println("Error reading page: ", cnt, " from file: ", pagesPath, err)
				return nil, err
			}
			if len(currentPage) == 0 {
				break
			}
			nodesInPage := 0
			for i := 1; i <= STORAGE_SUBHEAP_NODES; i++ {
				currentNode, err := deserializeNode(currentPage[(i-1)*STORAGE_NODE_SIZE_IN_BYTES : i*STORAGE_NODE_SIZE_IN_BYTES])
				if err != nil {
					fmt.Println("Error deserializing node: ", i, " from page: ", cnt, " in file: ", pagesPath, err)
					return nil, err
				}
				if currentNode.priority == 0 {
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

		return &KokaqCore{
			totalNodes:              nodes,
			totalPages:              pages,
			currentLoadedPageNumber: 0,
			currentLoadedPage:       nil,
			pagesPath:               pagesPath,
		}, nil
	}
}

func (queue *KokaqCore) Push(node *heapNode) error {
	if queue.totalNodes == 0 {
		queue.loadPage(0)
		newPage := make([]byte, STORAGE_SUBHEAP_SIZE_IN_BYTES)
		newNode, err := serializeNode(node)
		if err != nil {
			fmt.Println("Error serializing node", err)
			return err
		}
		copy(newPage[:STORAGE_NODE_SIZE_IN_BYTES], newNode)
		queue.saveNewPage(1, newPage)
	} else {
		queue.pageHeapifyUp(queue.totalNodes+1, node)
	}
	queue.totalNodes += 1
	return nil
}

func (queue *KokaqCore) Pop() (*heapNode, error) {
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

		newEle := make([]byte, STORAGE_NODE_SIZE_IN_BYTES)
		oldEle, err := deserializeNode(page[:STORAGE_NODE_SIZE_IN_BYTES])
		if err != nil {
			return nil, err
		}
		copy(page[:STORAGE_NODE_SIZE_IN_BYTES], newEle)
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
		startPoint := (localIndex - 1) * (STORAGE_NODE_SIZE_IN_BYTES)
		lastNode, err := deserializeNode(page[startPoint : startPoint+STORAGE_NODE_SIZE_IN_BYTES])
		if err != nil {
			return nil, err
		}
		newNode := make([]byte, STORAGE_NODE_SIZE_IN_BYTES)
		copy(page[startPoint:], newNode)
		queue.totalNodes -= 1
		if localIndex == 2 && pageNumber != 1 {
			queue.totalPages -= 1
			//Remove the root node as well
			copy(page[:STORAGE_NODE_SIZE_IN_BYTES], newNode)
		}
		err = queue.pageHeapifyDown(lastNode)
		if err != nil {
			return nil, err
		}
		return topNode, nil
	}
}

func (queue *KokaqCore) Peek() (*heapNode, error) {
	page, err := queue.loadPage(1)
	if err != nil {
		return nil, err
	}
	return deserializeNode(page[:STORAGE_NODE_SIZE_IN_BYTES])
}

func (queue *KokaqCore) SetIndexofPeekElement(index int) error {
	page, err := queue.loadPage(1)
	if err != nil {
		return err
	}
	node, err := deserializeNode(page[:STORAGE_NODE_SIZE_IN_BYTES])
	if err != nil {
		return err
	}

	node.indexPos = index
	newNode, err := serializeNode(node)
	if err != nil {
		return err
	}

	copy(page[:STORAGE_NODE_SIZE_IN_BYTES], newNode)
	return nil
}

func (queue *KokaqCore) pageHeapifyUp(index int, node *heapNode) error {
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
		rootNode := page[(localIndexOfParent-1)*(STORAGE_NODE_SIZE_IN_BYTES) : (localIndexOfParent)*(STORAGE_NODE_SIZE_IN_BYTES)]
		_, err = queue.loadPage(0)
		if err != nil {
			fmt.Println("Error loading page: 0 in file: ", queue.pagesPath, err)
			return err
		}
		newPage := make([]byte, STORAGE_SUBHEAP_SIZE_IN_BYTES)
		copy(newPage[:STORAGE_NODE_SIZE_IN_BYTES], rootNode)
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
		startPoint := (localIndex - 1) * (STORAGE_NODE_SIZE_IN_BYTES)
		nodeContent, err := serializeNode(node)
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
			utils.WriteBytesToFile(queue.pagesPath, int64((previousPage-1)*(STORAGE_SUBHEAP_SIZE_IN_BYTES)), page[startPoint:startPoint+STORAGE_NODE_SIZE_IN_BYTES])
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

func (queue *KokaqCore) localHeapifyUp(index int, heap []byte) (bool, error) {
	indexChild := index
	for indexChild > 1 {
		startPointChild := (indexChild - 1) * (STORAGE_NODE_SIZE_IN_BYTES)
		bytesChild := heap[startPointChild : startPointChild+STORAGE_NODE_SIZE_IN_BYTES]
		priorityChild := binary.LittleEndian.Uint64(bytesChild[:STORAGE_PRIORITY_SIZE_IN_BYTES])

		indexParent := indexChild / 2
		startPointParent := (indexParent - 1) * (STORAGE_NODE_SIZE_IN_BYTES)
		bytesParent := heap[startPointParent : startPointParent+STORAGE_NODE_SIZE_IN_BYTES]
		priorityParent := binary.LittleEndian.Uint64(bytesParent[:STORAGE_PRIORITY_SIZE_IN_BYTES])

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

func (queue *KokaqCore) pageHeapifyDown(node *heapNode) error {
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
		newNode, err := serializeNode(node)
		if err != nil {
			return err
		}
		copy(page[:STORAGE_NODE_SIZE_IN_BYTES], newNode)
		needToGoDown, lastLayerIndex, err = queue.localHeapifyDown(page)
		if err != nil {
			return err
		}
		//If second iteration change the last layer node of previous page
		//------Anti-pattern code------
		// Here we are directly going to previous page and changing the last layer node, every operation should be doen by loading page
		// Here, we are sure that page is synced in meomry, to reduce disk I/O this has been done
		if previousPage != 0 {
			utils.WriteBytesToFile(queue.pagesPath, int64((previousPage-1)*(STORAGE_SUBHEAP_SIZE_IN_BYTES)+(previousIndex-1)*(STORAGE_NODE_SIZE_IN_BYTES)), page[:STORAGE_NODE_SIZE_IN_BYTES])
		}

		//Change parametes for next iterations
		previousPage = pageNumber
		previousIndex = lastLayerIndex

		//Next page number which needs to be loaded based on current page we are on and last layer index
		pageNumber = (STORAGE_SUBHEAP_LAST_LAYER_NODES)*(pageNumber-1) + (lastLayerIndex - STORAGE_SUBHEAP_LAST_LAYER_NODES + 1) + 1

		//If no pages available, no need to go down
		if pageNumber > queue.totalPages {
			break
		}
	}
	return nil
}

func (queue *KokaqCore) localHeapifyDown(heap []byte) (bool, int, error) {
	indexParent := 1
	for indexParent < STORAGE_SUBHEAP_LAST_LAYER_NODES {
		parentBytes := heap[(indexParent-1)*(STORAGE_NODE_SIZE_IN_BYTES) : (indexParent)*(STORAGE_NODE_SIZE_IN_BYTES)]
		priorityParent := binary.LittleEndian.Uint64(parentBytes[:STORAGE_PRIORITY_SIZE_IN_BYTES])
		indexLeftChild := indexParent * 2
		indexRightChild := indexLeftChild + 1
		leftChildBytes := heap[(indexLeftChild-1)*(STORAGE_NODE_SIZE_IN_BYTES) : (indexLeftChild)*(STORAGE_NODE_SIZE_IN_BYTES)]
		rightChildBytes := heap[(indexRightChild-1)*(STORAGE_NODE_SIZE_IN_BYTES) : (indexRightChild)*(STORAGE_NODE_SIZE_IN_BYTES)]
		priorityLeftChild := binary.LittleEndian.Uint64(leftChildBytes[:STORAGE_PRIORITY_SIZE_IN_BYTES])
		priorityRightChild := binary.LittleEndian.Uint64(rightChildBytes[:STORAGE_PRIORITY_SIZE_IN_BYTES])

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
				copy(heap[(indexParent-1)*(STORAGE_NODE_SIZE_IN_BYTES):], leftChildBytes)
				copy(heap[(indexLeftChild-1)*(STORAGE_NODE_SIZE_IN_BYTES):], tempByteParents)
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
					copy(heap[(indexParent-1)*(STORAGE_NODE_SIZE_IN_BYTES):], leftChildBytes)
					copy(heap[(indexLeftChild-1)*(STORAGE_NODE_SIZE_IN_BYTES):], tempByteParents)
					indexParent = indexLeftChild
				} else {
					//Right child wins
					tempByteParents := make([]byte, len(parentBytes))
					copy(tempByteParents, parentBytes)
					copy(heap[(indexParent-1)*(STORAGE_NODE_SIZE_IN_BYTES):], rightChildBytes)
					copy(heap[(indexRightChild-1)*(STORAGE_NODE_SIZE_IN_BYTES):], tempByteParents)
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
func (queue *KokaqCore) loadPage(pageNumber int) ([]byte, error) {
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
		startIndex := (queue.currentLoadedPageNumber - 1) * (STORAGE_SUBHEAP_SIZE_IN_BYTES)
		err := utils.WriteBytesToFile(queue.pagesPath, int64(startIndex), queue.currentLoadedPage)
		if err != nil {
			fmt.Println("Error saving page: ", queue.currentLoadedPageNumber, " to file: ", queue.pagesPath, err)
			return nil, err
		}
	}
	if pageNumber != 0 {
		startIndex := (pageNumber - 1) * (STORAGE_SUBHEAP_SIZE_IN_BYTES)
		var err error
		queue.currentLoadedPage, err = utils.ReadBytesFromFile(queue.pagesPath, startIndex, STORAGE_SUBHEAP_SIZE_IN_BYTES)
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
func (queue *KokaqCore) saveNewPage(pageNumber int, data []byte) error {
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
func (queue *KokaqCore) getLocalHeapDetailsForNode(index int) (int, int, int) {
	globalLevel := bits.Len(uint(index)) - 1
	nodesInGloablLevel := 1 << globalLevel
	subHeapLevel := (globalLevel - 1) / (STORAGE_HEAP_MAX_SIZE - 1)
	localLevel := ((globalLevel - 1) % (STORAGE_HEAP_MAX_SIZE - 1)) + 1
	nodesInLocalLevel := 1 << localLevel
	pPartial := (index - nodesInGloablLevel) / (nodesInLocalLevel)
	pFull := (utils.Power(STORAGE_SUBHEAP_LAST_LAYER_NODES, subHeapLevel) - 1) / (STORAGE_SUBHEAP_LAST_LAYER_NODES - 1)
	p := pFull + pPartial + 1
	localIndexFull := nodesInLocalLevel - 1
	localIndexPartial := (index - nodesInGloablLevel) % nodesInLocalLevel
	localIndex := localIndexFull + localIndexPartial + 1
	return p, localIndex, localLevel
}

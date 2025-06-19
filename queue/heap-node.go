package queue

import (
	"encoding/binary"
)

type HeapNode struct {
	Priority int // Priority of the node
	Index    int // Position of the node in the heap
}

func deserializeNode(data []byte, prioritySizeInBytes int) (*HeapNode, error) {
	priority := binary.LittleEndian.Uint64(data[:prioritySizeInBytes])
	indexPos := binary.LittleEndian.Uint64(data[prioritySizeInBytes:])
	return NewHeapNode(int(priority), int(indexPos)), nil
}

func serializeNode(node *HeapNode, prioritySizeInBytes int, indexSizeInBytes int, nodeSizeInBytes int) ([]byte, error) {
	seralizedNode := make([]byte, nodeSizeInBytes)

	buf := make([]byte, prioritySizeInBytes)
	binary.LittleEndian.PutUint64(buf, uint64(node.Priority))
	copy(seralizedNode[0:], buf)

	buf = make([]byte, indexSizeInBytes)
	binary.LittleEndian.PutUint64(buf, uint64(node.Index))
	copy(seralizedNode[prioritySizeInBytes:], buf)
	return seralizedNode, nil
}

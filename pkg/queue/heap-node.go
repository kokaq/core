package queue

import (
	"encoding/binary"
)

func deserializeNode(data []byte) (*heapNode, error) {
	priority := binary.LittleEndian.Uint64(data[:STORAGE_PRIORITY_SIZE_IN_BYTES])
	indexPos := binary.LittleEndian.Uint64(data[STORAGE_PRIORITY_SIZE_IN_BYTES:])
	return &heapNode{
		priority: int(priority),
		indexPos: int(indexPos),
	}, nil
}

func serializeNode(node *heapNode) ([]byte, error) {
	seralizedNode := make([]byte, STORAGE_NODE_SIZE_IN_BYTES)

	buf := make([]byte, STORAGE_PRIORITY_SIZE_IN_BYTES)
	binary.LittleEndian.PutUint64(buf, uint64(node.priority))
	copy(seralizedNode[0:], buf)

	buf = make([]byte, STORAGE_INDEX_SIZE_IN_BYTES)
	binary.LittleEndian.PutUint64(buf, uint64(node.indexPos))
	copy(seralizedNode[STORAGE_PRIORITY_SIZE_IN_BYTES:], buf)
	return seralizedNode, nil
}

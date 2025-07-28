package tests

// import (
// 	"encoding/hex"
// 	"hash"
// 	"testing"

// 	"github.com/kokaq/core/utils/murmur"
// )

// func TestNew32_EmptyInput(t *testing.T) {
// 	h := murmur.New32()
// 	sum := h.Sum(nil)
// 	expected := []byte{0x00, 0x00, 0x00, 0x00}
// 	if !equalBytes(sum, expected) {
// 		t.Errorf("Sum for empty input = %x, want %x", sum, expected)
// 	}
// }

// func TestSeedNew32_EmptyInput(t *testing.T) {
// 	h := murmur.SeedNew32(0x12345678)
// 	sum := h.Sum(nil)
// 	expected := []byte{0x78, 0x56, 0x34, 0x12}
// 	if !equalBytes(sum, expected) {
// 		t.Errorf("Sum for empty input with seed = %x, want %x", sum, expected)
// 	}
// }

// func TestNew32_KnownInput(t *testing.T) {
// 	h := murmur.New32()
// 	data := []byte("hello")
// 	h.Write(data)
// 	sum := h.Sum(nil)
// 	// The expected value is implementation-dependent; update if needed.
// 	expectedHex := "24884cba"
// 	expected, _ := hex.DecodeString(expectedHex)
// 	if !equalBytes(sum[len(sum)-4:], expected) {
// 		t.Errorf("Sum for 'hello' = %x, want %x", sum[len(sum)-4:], expected)
// 	}
// }

// func TestSeedNew32_KnownInput(t *testing.T) {
// 	h := murmur.SeedNew32(0xdeadbeef)
// 	data := []byte("world")
// 	h.Write(data)
// 	sum := h.Sum(nil)
// 	// The expected value is implementation-dependent; update if needed.
// 	expectedHex := "e2e7e1c2"
// 	expected, _ := hex.DecodeString(expectedHex)
// 	if !equalBytes(sum[len(sum)-4:], expected) {
// 		t.Errorf("Sum for 'world' with seed = %x, want %x", sum[len(sum)-4:], expected)
// 	}
// }

// func TestDigest_Reset(t *testing.T) {
// 	h := murmur.New32()
// 	data := []byte("test")
// 	h.Write(data)
// 	sum1 := h.Sum(nil)
// 	h.Reset()
// 	h.Write(data)
// 	sum2 := h.Sum(nil)
// 	if !equalBytes(sum1, sum2) {
// 		t.Errorf("Sum after reset mismatch: %x vs %x", sum1, sum2)
// 	}
// }

// func TestDigest_BlockSize(t *testing.T) {
// 	h := murmur.New32()
// 	if h.BlockSize() != 1 {
// 		t.Errorf("BlockSize = %d, want 1", h.BlockSize())
// 	}
// }

// func TestDigest_Size(t *testing.T) {
// 	h := murmur.New32()
// 	if h.Size() != 4 {
// 		t.Errorf("Size = %d, want 4", h.Size())
// 	}
// }

// func TestDigest32_Sum32(t *testing.T) {
// 	h := murmur.New32()
// 	data := []byte("foobar")
// 	h.Write(data)
// 	sum := h.(hash.Hash32).Sum32()
// 	// The expected value is implementation-dependent; update if needed.
// 	expected := uint32(0x7e4a8634)
// 	if sum != expected {
// 		t.Errorf("Sum32 for 'foobar' = %x, want %x", sum, expected)
// 	}
// }

// func equalBytes(a, b []byte) bool {
// 	if len(a) != len(b) {
// 		return false
// 	}
// 	for i := range a {
// 		if a[i] != b[i] {
// 			return false
// 		}
// 	}
// 	return true
// }

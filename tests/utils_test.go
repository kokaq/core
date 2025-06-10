package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kokaq/core/v1/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestPower(t *testing.T) {
	assert.Equal(t, 1, utils.Power(2, 0))
	assert.Equal(t, 2, utils.Power(2, 1))
	assert.Equal(t, 8, utils.Power(2, 3))
	assert.Equal(t, 27, utils.Power(3, 3))
	assert.Equal(t, 1, utils.Power(1, 100))
	assert.Equal(t, 0, utils.Power(0, 5))
}

func TestLog2(t *testing.T) {
	assert.Equal(t, -1, utils.Log2(0))
	assert.Equal(t, 0, utils.Log2(1))
	assert.Equal(t, 1, utils.Log2(2))
	assert.Equal(t, 2, utils.Log2(4))
	assert.Equal(t, 3, utils.Log2(8))
	assert.Equal(t, 4, utils.Log2(16))
	assert.Equal(t, 5, utils.Log2(31)) // 2^5 > 31
	assert.Equal(t, 5, utils.Log2(32))
}

func TestIsPowerOfTwo(t *testing.T) {
	assert.False(t, utils.IsPowerOfTwo(0))
	assert.True(t, utils.IsPowerOfTwo(1))
	assert.True(t, utils.IsPowerOfTwo(2))
	assert.False(t, utils.IsPowerOfTwo(3))
	assert.True(t, utils.IsPowerOfTwo(4))
	assert.False(t, utils.IsPowerOfTwo(6))
	assert.True(t, utils.IsPowerOfTwo(8))
}

// --- Fuzz Tests (for mathematical functions) ---

func FuzzPower(f *testing.F) {
	f.Add(2, 3)
	f.Add(5, 0)
	f.Add(1, 100)
	f.Fuzz(func(t *testing.T, base, exp int) {
		if exp < 0 {
			t.Skip() // Skip negative exponents
		}
		result := utils.Power(base, exp)
		assert.GreaterOrEqual(t, result, 0)
	})
}

func FuzzLog2(f *testing.F) {
	f.Add(1)
	f.Add(2)
	f.Add(1024)
	f.Fuzz(func(t *testing.T, x int) {
		if x <= 0 {
			t.Skip()
		}
		logVal := utils.Log2(x)
		assert.True(t, 1<<logVal <= x)
	})
}

func FuzzIsPowerOfTwo(f *testing.F) {
	f.Add(1)
	f.Add(2)
	f.Add(3)
	f.Add(1024)
	f.Fuzz(func(t *testing.T, x int) {
		isPow := utils.IsPowerOfTwo(x)
		if isPow {
			assert.Equal(t, 0, x&(x-1))
		}
	})
}

func TestFileOperations(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "testfile.bin")

	// Test WriteBytesToFile
	data := []byte("hello world")
	err := utils.WriteBytesToFile(filePath, 0, data)
	assert.NoError(t, err)

	// Test ReadBytesFromFile
	readData, err := utils.ReadBytesFromFile(filePath, 0, len(data))
	assert.NoError(t, err)
	assert.Equal(t, data, readData)

	// Test ReadBytesFromFile with offset
	offsetData := []byte("world")
	readData, err = utils.ReadBytesFromFile(filePath, 6, 5)
	assert.NoError(t, err)
	assert.Equal(t, offsetData, readData)

	// Test AppendBytesToFile
	appendData := []byte("!!!")
	err = utils.AppendBytesToFile(filePath, appendData)
	assert.NoError(t, err)

	expected := append(data, appendData...)
	readData, err = os.ReadFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, expected, readData)

	// Test WriteBytesToFile with offset (overwrite)
	overwriteData := []byte("HELLO")
	err = utils.WriteBytesToFile(filePath, 0, overwriteData)
	assert.NoError(t, err)

	readData, err = utils.ReadBytesFromFile(filePath, 0, len(overwriteData))
	assert.NoError(t, err)
	assert.Equal(t, overwriteData, readData)
}

func TestReadBytesFromFile_Errors(t *testing.T) {
	// Non-existent file
	_, err := utils.ReadBytesFromFile("nonexistent.txt", 0, 10)
	assert.Error(t, err)

	// Create a file and close it to test seek/read errors
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "testfile_error.bin")
	err = os.WriteFile(filePath, []byte("short"), 0644)
	assert.NoError(t, err)

	// Offset beyond file length
	_, err = utils.ReadBytesFromFile(filePath, 100, 10)
	assert.NoError(t, err) // Seeking past EOF is allowed, reading returns 0 bytes
}

func TestWriteBytesToFile_Errors(t *testing.T) {
	// Attempt to write to a directory
	tmpDir := t.TempDir()
	err := utils.WriteBytesToFile(tmpDir, 0, []byte("invalid"))
	assert.Error(t, err)
}

func TestAppendBytesToFile_Errors(t *testing.T) {
	// Attempt to append to a directory
	tmpDir := t.TempDir()
	err := utils.AppendBytesToFile(tmpDir, []byte("invalid"))
	assert.Error(t, err)
}

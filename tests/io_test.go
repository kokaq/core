package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kokaq/core/utils"
	"github.com/stretchr/testify/assert"
)

func TestDirectoryExists(t *testing.T) {
	tmpDir := t.TempDir()
	assert.True(t, utils.DirectoryExists(tmpDir))
	assert.False(t, utils.DirectoryExists(filepath.Join(tmpDir, "nonexistent")))
}

func TestEnsureDirectoryCreated(t *testing.T) {
	tmpDir := t.TempDir()
	newDir := filepath.Join(tmpDir, "subdir")
	err := utils.EnsureDirectoryCreated(newDir)
	assert.NoError(t, err)
	assert.True(t, utils.DirectoryExists(newDir))

	// Should not error if already exists
	err = utils.EnsureDirectoryCreated(newDir)
	assert.NoError(t, err)
}

func TestEnsureDirectoryDeleted(t *testing.T) {
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "subdir")
	os.Mkdir(subDir, 0755)
	assert.True(t, utils.DirectoryExists(subDir))

	err := utils.EnsureDirectoryDeleted(subDir)
	assert.NoError(t, err)
	assert.False(t, utils.DirectoryExists(subDir))

	// Should not error if already deleted
	err = utils.EnsureDirectoryDeleted(subDir)
	assert.NoError(t, err)
}

// func TestFileExists(t *testing.T) {
// 	tmpDir := t.TempDir()
// 	filePath := filepath.Join(tmpDir, "file.txt")
// 	assert.False(t, utils.FileExists(filePath))

// 	err := ioutil.WriteFile(filePath, []byte("data"), 0644)
// 	assert.NoError(t, err)
// 	assert.True(t, utils.FileExists(filePath))

// 	// Directory should not be considered a file
// 	assert.False(t, utils.FileExists(tmpDir))
// }

func TestEnsureFileCreated(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "file.txt")
	err := utils.EnsureFileCreated(filePath)
	assert.NoError(t, err)
	assert.True(t, utils.FileExists(filePath))

	// Should not error if already exists
	err = utils.EnsureFileCreated(filePath)
	assert.NoError(t, err)
}

// func TestEnsureFileDeleted(t *testing.T) {
// 	tmpDir := t.TempDir()
// 	filePath := filepath.Join(tmpDir, "file.txt")
// 	ioutil.WriteFile(filePath, []byte("data"), 0644)
// 	assert.True(t, utils.FileExists(filePath))

// 	err := utils.EnsureFileDeleted(filePath)
// 	assert.NoError(t, err)
// 	assert.False(t, utils.FileExists(filePath))

// 	// Should not error if already deleted
// 	err = utils.EnsureFileDeleted(filePath)
// 	assert.NoError(t, err)
// }

// func TestReadBytesFromFile(t *testing.T) {
// 	tmpDir := t.TempDir()
// 	filePath := filepath.Join(tmpDir, "file.txt")
// 	content := []byte("hello world")
// 	ioutil.WriteFile(filePath, content, 0644)

// 	// Read first 5 bytes
// 	data, err := utils.ReadBytesFromFile(filePath, 0, 5)
// 	assert.NoError(t, err)
// 	assert.Equal(t, []byte("hello"), data)

// 	// Read with offset
// 	data, err = utils.ReadBytesFromFile(filePath, 6, 5)
// 	assert.NoError(t, err)
// 	assert.Equal(t, []byte("world"), data)

// 	// Read beyond EOF
// 	data, err = utils.ReadBytesFromFile(filePath, int64(len(content)), 5)
// 	assert.NoError(t, err)
// 	assert.Equal(t, 0, len(data))

// 	// Read from non-existent file
// 	_, err = utils.ReadBytesFromFile(filepath.Join(tmpDir, "nofile.txt"), 0, 5)
// 	assert.Error(t, err)
// }

// func TestWriteBytesToFile(t *testing.T) {
// 	tmpDir := t.TempDir()
// 	filePath := filepath.Join(tmpDir, "file.txt")
// 	data := []byte("abcdef")

// 	// Write at offset 0
// 	err := utils.WriteBytesToFile(filePath, 0, data)
// 	assert.NoError(t, err)
// 	read, _ := ioutil.ReadFile(filePath)
// 	assert.Equal(t, data, read)

// 	// Overwrite at offset 3
// 	err = utils.WriteBytesToFile(filePath, 3, []byte("XYZ"))
// 	assert.NoError(t, err)
// 	read, _ = ioutil.ReadFile(filePath)
// 	assert.Equal(t, []byte("abcXYZ"), read)
// }

// func TestAppendBytesToFile(t *testing.T) {
// 	tmpDir := t.TempDir()
// 	filePath := filepath.Join(tmpDir, "file.txt")

// 	// Append to new file
// 	err := utils.AppendBytesToFile(filePath, []byte("foo"))
// 	assert.NoError(t, err)
// 	read, _ := ioutil.ReadFile(filePath)
// 	assert.Equal(t, []byte("foo"), read)

// 	// Append again
// 	err = utils.AppendBytesToFile(filePath, []byte("bar"))
// 	assert.NoError(t, err)
// 	read, _ = ioutil.ReadFile(filePath)
// 	assert.Equal(t, []byte("foobar"), read)
// }

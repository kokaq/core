package utils

import (
	"fmt"
	"io"
	"os"
)

func DirectoryExists(path string) bool {
	// This function checks if a directory exists at the given path.
	// Returns true if it exists, false otherwise.
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func EnsureDirectoryCreated(path string) error {
	// This function ensures that the directory at the given path exists.
	// If it does not exist, it creates the directory.
	// If it exists, it does nothing.
	// Returns an error if the directory cannot be created or accessed.

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory %s: %w", path, err)
		}
	}
	return nil
}

func EnsureDirectoryDeleted(path string) error {
	// This function ensures that the directory at the given path is deleted.
	// If it does not exist, it does nothing.
	// Returns an error if the directory cannot be deleted.

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		err = os.RemoveAll(path)
		if err != nil {
			return fmt.Errorf("failed to delete directory %s: %w", path, err)
		}
	}
	return nil
}

func FileExists(path string) bool {
	// This function checks if a file exists at the given path.
	// Returns true if it exists, false otherwise.
	info, err := os.Stat(path)
	return !os.IsNotExist(err) && !info.IsDir()
}

func EnsureFileCreated(path string) error {
	// This function ensures that the file at the given path exists.
	// If it does not exist, it creates an empty file.
	// Returns an error if the file cannot be created or accessed.

	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", path, err)
		}
		defer file.Close()
	}
	return nil
}

func EnsureFileDeleted(path string) error {
	// This function ensures that the file at the given path is deleted.
	// If it does not exist, it does nothing.
	// Returns an error if the file cannot be deleted.

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		err = os.Remove(path)
		if err != nil {
			return fmt.Errorf("failed to delete file %s: %w", path, err)
		}
	}
	return nil
}

func ReadBytesFromFile(path string, offset int64, length int) ([]byte, error) {
	// This function reads a specified number of bytes from a file at a given offset.
	// Returns the read bytes or an error if the operation fails.

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", path, err)
	}
	defer file.Close()

	_, err = file.Seek(int64(offset), io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("failed to seek in file %s: %w", path, err)
	}

	buffer := make([]byte, length)
	bytesRead, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("failed to read bytes from file %s: %w", path, err)
	}

	// return the bytes read, trimming any excess if less than requested length
	return buffer[:bytesRead], nil
}

func WriteBytesToFile(path string, offset int64, data []byte) error {
	// This function writes a byte slice to a file at a specified offset.
	// If the file does not exist, it creates it.
	// Returns an error if the operation fails.

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", path, err)
	}
	defer file.Close()

	_, err = file.Seek(int64(offset), io.SeekStart)
	if err != nil {
		return fmt.Errorf("failed to seek at %d in file %s: %w", offset, path, err)
	}
	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write bytes to file %s: %w", path, err)
	}
	return nil
}

func AppendBytesToFile(path string, data []byte) error {
	// This function appends a byte slice to a file.
	// If the file does not exist, it creates it.
	// Returns an error if the operation fails.

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", path, err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("failed to append bytes to file %s: %w", path, err)
	}
	return nil
}

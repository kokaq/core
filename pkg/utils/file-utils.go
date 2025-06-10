package utils

import (
	"fmt"
	"io"
	"os"
)

func ReadBytesFromFile(filePath string, offset int, length int) ([]byte, error) {
	// open the file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file: ", filePath, err)
		return nil, err
	}
	defer file.Close()

	// seek to the starting position
	_, err = file.Seek(int64(offset), io.SeekStart)
	if err != nil {
		fmt.Println("Error seeking bytes at: ", offset, " in file: ", filePath, err)
		return nil, err
	}

	// read the specified number of bytes
	buffer := make([]byte, length)
	bytesRead, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		fmt.Println("Error reading bytes at: ", offset, " in file: ", filePath, err)
		return nil, err
	}

	// return the bytes read, trimming any excess if less than requested length
	return buffer[:bytesRead], nil
}

func WriteBytesToFile(filePath string, offset int64, data []byte) error {
	// Open the file with read/write permissions
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening file: ", filePath, err)
		return err
	}
	defer file.Close()

	// Seek to the specified offset
	_, err = file.Seek(offset, io.SeekStart)
	if err != nil {
		fmt.Println("Error seeking bytes at: ", offset, " in file: ", filePath, err)
		return err
	}

	// Write the bytes to the file at the current offset
	_, err = file.Write(data)
	if err != nil {
		fmt.Println("Error writing bytes at: ", offset, " in file: ", filePath, err)
		return err
	}
	return nil
}

func AppendBytesToFile(filePath string, data []byte) error {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", filePath, err)
		return err
	}
	defer f.Close()

	if _, err := f.Write(data); err != nil {
		fmt.Println("Error appending to file: ", filePath, err)
		return err
	}
	return nil
}

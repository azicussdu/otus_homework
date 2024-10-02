package main

import (
	"errors"
	"fmt"
	"github.com/cheggaaa/pb/v3" //nolint:depguard
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func closeFile(file *os.File) {
	err := file.Close()
	if err != nil {
		_ = fmt.Errorf("failed to close file: %w", err)
	}
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	readFile, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer closeFile(readFile)

	fileInfo, _ := readFile.Stat()
	fileSize := fileInfo.Size()

	if !fileInfo.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}

	writeFile, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer closeFile(writeFile)

	// move the file pointer to offset bytes from beginning of file
	if _, err = readFile.Seek(offset, io.SeekStart); err != nil {
		return fmt.Errorf("failed to seek in file: %w", err)
	}

	if limit == 0 || limit > (fileSize-offset) {
		limit = fileSize - offset
	}

	bar := pb.Full.Start64(limit)
	barReader := bar.NewProxyReader(readFile)

	// buffer size is 10KB so it can copy the whole content of input.txt file
	// in case if fromPath == toPath
	bufferSize := 10 * 1024
	buffer := make([]byte, bufferSize)

	totalCopied := int64(0)

	for totalCopied < limit {
		bytesRemaining := limit - totalCopied

		bytesToRead := bufferSize
		if bytesRemaining < int64(bufferSize) {
			bytesToRead = int(bytesRemaining)
		}

		var n int
		n, err = barReader.Read(buffer[:bytesToRead]) // read only bytesToRead bytes
		if err != nil && err != io.EOF {
			return fmt.Errorf("failed to read data: %w", err)
		}
		if n == 0 {
			break // If data to read is finished
		}

		_, err = writeFile.Write(buffer[:n]) // write only n bytes
		if err != nil {
			err = os.Remove(writeFile.Name()) // remove file if write fails
			if err != nil {
				return err
			}
			return fmt.Errorf("failed to write data: %w", err)
		}

		// Update the total number of bytes copied
		totalCopied += int64(n)
	}
	bar.Finish()

	return nil
}

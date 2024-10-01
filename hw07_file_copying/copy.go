package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/cheggaaa/pb/v3" //nolint:depguard
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

	if filepath.Ext(fromPath) != ".txt" {
		return ErrUnsupportedFile
	}

	fileInfo, _ := readFile.Stat()
	fileSize := fileInfo.Size()

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

	bytesToRead := fileSize - offset
	// As I got it if limit is 0 we have to copy the whole file
	if limit == 0 || limit > bytesToRead {
		limit = bytesToRead
	}

	bar := pb.Full.Start64(limit)
	barReader := bar.NewProxyReader(readFile)

	_, err = io.CopyN(writeFile, barReader, limit)
	if err != nil && errors.Is(err, io.EOF) {
		return fmt.Errorf("failed to copy data: %w", err)
	}

	return nil
}

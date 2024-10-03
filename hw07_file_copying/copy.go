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

func copyContent(readFile, writeFile *os.File, offset, limit, fileSize int64) error {
	if _, err := readFile.Seek(offset, io.SeekStart); err != nil {
		return fmt.Errorf("failed to seek in file: %w", err)
	}

	if limit == 0 || limit > (fileSize-offset) {
		limit = fileSize - offset
	}

	bar := pb.Full.Start64(limit)
	barReader := bar.NewProxyReader(readFile)

	_, err := io.CopyN(writeFile, barReader, limit)
	if err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("failed to copy data: %w", err)
	}

	bar.Finish()
	return nil
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

	if fromPath == toPath {
		readFile, err = handleTempFileCopy(toPath, readFile, offset, limit, fileSize)
		if err != nil {
			return err
		}
		offset = 0
	}

	writeFile, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer closeFile(writeFile)

	if err = copyContent(readFile, writeFile, offset, limit, fileSize); err != nil {
		return err
	}

	return nil
}

func handleTempFileCopy(toPath string, readFile *os.File, offset, limit, fileSize int64) (*os.File, error) {
	tmpFile, err := os.CreateTemp(filepath.Dir(toPath), "tempfile-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer closeFile(tmpFile)

	if err = copyContent(readFile, tmpFile, offset, limit, fileSize); err != nil {
		return nil, err
	}
	if err = tmpFile.Close(); err != nil {
		return nil, fmt.Errorf("failed to close temporary file: %w", err)
	}

	tmpFile, err = os.Open(tmpFile.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to reopen temporary file: %w", err)
	}

	defer func() {
		err = os.Remove(tmpFile.Name())
		if err != nil {
			return
		}
	}()

	return tmpFile, nil
}

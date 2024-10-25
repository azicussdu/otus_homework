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
	absFromPath, err := filepath.Abs(fromPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for fromPath: %w", err)
	}

	absToPath, err := filepath.Abs(toPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for toPath: %w", err)
	}

	readFile, err := os.Open(absFromPath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}

	fileInfo, _ := readFile.Stat()
	fileSize := fileInfo.Size()

	if !fileInfo.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}

	tmpFile, err := os.CreateTemp(filepath.Dir(absToPath), "tempfile-*")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}

	if err = copyContent(readFile, tmpFile, offset, limit, fileSize); err != nil {
		return err
	}

	if err = tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temporary file: %w", err)
	}

	if err = readFile.Close(); err != nil {
		return fmt.Errorf("failed to close the original file: %w", err)
	}

	if absFromPath == absToPath {
		if err = os.Remove(absFromPath); err != nil {
			return fmt.Errorf("failed to remove original file: %w", err)
		}
	}

	if err = os.Rename(tmpFile.Name(), absToPath); err != nil {
		return fmt.Errorf("failed to rename temporary file: %w", err)
	}

	return nil
}

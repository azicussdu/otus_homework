package main

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require" //nolint:depguard
)

func TestCopy(t *testing.T) {
	fromPath := "testdata/input.txt"
	toPath := "testdata/output.txt"

	offsets := []int64{0, 0, 0, 0, 100, 6000}
	limits := []int64{0, 10, 1000, 10000, 1000, 1000}

	for i := 0; i < 6; i++ {
		toPathCheck := fmt.Sprintf("testdata/out_offset%d_limit%d.txt", offsets[i], limits[i])

		err := Copy(fromPath, toPath, offsets[i], limits[i])
		if err != nil {
			t.Fatalf("Failed to copy file = %v", err)
		}

		file1Info, err := os.Stat(toPath)
		if err != nil {
			t.Fatalf("Failed to get file info = %v", err)
		}
		file1Size := file1Info.Size()

		file2Info, err := os.Stat(toPathCheck)
		if err != nil {
			t.Fatalf("Failed to get file info = %v", err)
		}
		file2Size := file2Info.Size()

		require.Equal(t, file1Size, file2Size)

		if file1Size >= 10 {
			same, _ := compareStartEnd10Bytes(toPath, toPathCheck)
			require.True(t, same)
		}
	}

	err := os.Remove(toPath)
	if err != nil {
		return
	}
}

func compareStartEnd10Bytes(fromFile, toFile string) (bool, error) {
	from, err := os.Open(fromFile)
	if err != nil {
		return false, err
	}
	defer from.Close()

	to, err := os.Open(toFile)
	if err != nil {
		return false, err
	}
	defer to.Close()

	fromStart := make([]byte, 10)
	_, err = from.Read(fromStart)
	if err != nil {
		return false, err
	}

	toStart := make([]byte, 10)
	_, err = to.Read(toStart)
	if err != nil {
		return false, err
	}

	if string(fromStart) != string(toStart) {
		return false, errors.New("Start 10 bytes are different")
	}

	fromEnd := make([]byte, 10)
	_, err = from.Seek(-10, 2)
	if err != nil {
		return false, err
	}
	_, err = from.Read(fromEnd)
	if err != nil {
		return false, err
	}

	toEnd := make([]byte, 10)
	_, err = to.Seek(-10, 2)
	if err != nil {
		return false, err
	}
	_, err = to.Read(toEnd)
	if err != nil {
		return false, err
	}

	if string(fromEnd) != string(toEnd) {
		return false, errors.New("Start 10 bytes are different")
	}

	return true, nil
}

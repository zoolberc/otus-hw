package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("Copy file with offset 0 and limit 16", func(t *testing.T) {
		err := Copy("copytestdata/input_file.txt", "result1.txt", 0, 16)
		require.NoError(t, err)
		tByte, errReadFile := os.ReadFile("result1.txt")
		require.NoError(t, errReadFile)
		require.Equal(t, "IT professionals", bytes.NewBuffer(tByte).String())
		errRemoveFile := os.Remove("result1.txt")
		require.NoError(t, errRemoveFile)
	})

	t.Run("Copy file with offset 3 and limit 17", func(t *testing.T) {
		err := Copy("copytestdata/input_file.txt", "result2.txt", 3, 17)
		require.NoError(t, err)
		tByte, errReadFile := os.ReadFile("result2.txt")
		require.NoError(t, errReadFile)
		require.Equal(t, "professionals are", bytes.NewBuffer(tByte).String())
		errRemoveFile := os.Remove("result2.txt")
		require.NoError(t, errRemoveFile)
	})

	t.Run("Copy file with offset 0 and limit 0", func(t *testing.T) {
		err := Copy("copytestdata/input_file.txt", "result3.txt", 0, 0)
		require.NoError(t, err)
		tByteNewFile, errReadNewFile := os.ReadFile("result3.txt")
		tByteTargetFile, errReadTargetFile := os.ReadFile("copytestdata/input_file.txt")
		require.NoError(t, errReadNewFile)
		require.NoError(t, errReadTargetFile)
		require.Equal(t, tByteTargetFile, tByteNewFile)
		errRemoveFile := os.Remove("result3.txt")
		require.NoError(t, errRemoveFile)
	})

	t.Run("Get error 'not found file'", func(t *testing.T) {
		err := Copy("copytestdata/error.txt", "result4.txt", 0, 0)
		require.ErrorIs(t, err, ErrFileNotFound)
	})

	t.Run("Get error 'offset exceeds file size'", func(t *testing.T) {
		err := Copy("copytestdata/input_file.txt", "result5.txt", 1000, 10)
		require.ErrorIs(t, err, ErrOffsetExceedsFileSize)
	})
}

package main

import (
	"errors"
	"io"
	"os"
	"strings"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrFileNotFound          = errors.New("file not found")
	ErrCreateNewFile         = errors.New("error when creating a new file")
	ErrOffsetLessZero        = errors.New("offset can`t be less than 0")
	ErrLimitLessZero         = errors.New("limit can`t less than 0")
	ErrTargetFileNotFound    = errors.New("target file not found")
	ErrNewFileNotFound       = errors.New("new file name can`t be empty")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if offset < 0 {
		return ErrOffsetLessZero
	}
	if limit < 0 {
		return ErrLimitLessZero
	}
	if strings.TrimSpace(fromPath) == "" {
		return ErrTargetFileNotFound
	}
	if strings.TrimSpace(toPath) == "" {
		return ErrNewFileNotFound
	}

	targetFile, err := os.Open(fromPath)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrFileNotFound
		}
		return err
	}
	fileStat, err := targetFile.Stat()
	if err != nil {
		return err
	}

	if offset > fileStat.Size() {
		return ErrOffsetExceedsFileSize
	}

	if offset > 0 {
		if _, err = targetFile.Seek(offset, io.SeekStart); err != nil {
			return err
		}
	}
	newFile, err := os.Create(toPath)
	if err != nil {
		return ErrCreateNewFile
	}

	if limit == 0 {
		limit = fileStat.Size()
	}

	reader := io.LimitReader(targetFile, limit)
	bar := pb.Full.Start64(limit)
	barReader := bar.NewProxyReader(reader)
	if _, err = io.CopyN(newFile, barReader, limit); err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}
	defer bar.Finish()
	defer targetFile.Close()
	defer newFile.Close()

	return nil
}

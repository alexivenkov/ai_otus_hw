package main

import (
	"errors"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrInvalidOffset         = errors.New("offset should be positive integer")
	ErrInvalidLimit          = errors.New("limit should be positive integer")
)

func Copy(fromPath, toPath string, offset, limit int64, disableOutput bool) error {
	if err := validateInput(offset, limit); err != nil {
		return err
	}

	src, err := os.Open(fromPath)
	if err != nil {
		return err
	}

	stat, err := src.Stat()
	if err != nil || !stat.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	dst, err := os.Create(toPath)
	if err != nil {
		return err
	}

	defer func() {
		src.Close()
		dst.Close()
	}()

	if stat.Size() < offset {
		return ErrOffsetExceedsFileSize
	}

	if _, err := src.Seek(offset, io.SeekStart); err != nil {
		return err
	}

	if limit == 0 || limit > stat.Size() {
		limit = stat.Size()
	}

	var (
		current    int64
		bufferSize int64 = 1000
	)

	if limit < bufferSize {
		bufferSize = limit
	}

	bar := new(Bar)
	bar.Init(limit, disableOutput)

	for current < limit {
		bytesCopied, err := io.CopyN(dst, src, bufferSize)
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return err
		}

		current += bytesCopied
		if err := bar.Progress(bytesCopied); err != nil {
			return err
		}
		bar.Print()
	}
	bar.Finish()

	return nil
}

func validateInput(offset, limit int64) error {
	if offset < 0 {
		return ErrInvalidOffset
	}
	if limit < 0 {
		return ErrInvalidLimit
	}
	return nil
}

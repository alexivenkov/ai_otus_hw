package main

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	src, _ := os.Open("testdata/input.txt")
	content, _ := ioutil.ReadFile("testdata/input.txt")
	stat, _ := src.Stat()
	src.Close()

	offset := content[1000:]

	testCases := []struct {
		testName        string
		offset          int64
		limit           int64
		expectedSize    int64
		expectedContent []byte
	}{
		{
			testName:        "Copy whole file",
			offset:          0,
			limit:           0,
			expectedSize:    stat.Size(),
			expectedContent: content,
		},
		{
			testName:        "Limit 1000",
			offset:          0,
			limit:           1000,
			expectedSize:    1000,
			expectedContent: content[:1000],
		},
		{
			testName:        "Offset 1000",
			offset:          1000,
			limit:           0,
			expectedSize:    stat.Size() - 1000,
			expectedContent: content[1000:],
		},
		{
			testName:        "Offset and limit 1000",
			offset:          1000,
			limit:           1000,
			expectedSize:    1000,
			expectedContent: offset[:1000],
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.testName, func(t *testing.T) {
			err := Copy("testdata/input.txt", "testdata/copied.txt", tc.offset, tc.limit, true)
			require.NoError(t, err)

			copiedFile, _ := os.Open("testdata/copied.txt")
			content, _ := ioutil.ReadFile("testdata/copied.txt")
			copiedFileStat, _ := copiedFile.Stat()
			copiedFile.Close()
			require.Equal(t, tc.expectedSize, copiedFileStat.Size())
			require.Equal(t, tc.expectedContent, content)
		})
	}
	defer os.Remove("testdata/copied.txt")
}

func TestCannotCopyDir(t *testing.T) {
	defer func() {
		os.Remove("testdata/testdir")
	}()
	os.Mkdir("testdata/testdir", 0777)

	err := Copy("testdata/testdir", "testdata/testdir_copy", 0, 1000, true)
	require.Truef(t, errors.Is(err, ErrUnsupportedFile), "actual error %q", err)
}

func TestUnknownFileLength(t *testing.T) {
	err := Copy("/dev/urandom", "testdata/test_copy", 0, 1000, true)
	require.Truef(t, errors.Is(err, ErrUnsupportedFile), "actual error %q", err)
}

func TestWrongPermissions(t *testing.T) {
	f, _ := os.Create("testdata/test_wrong_permissions")
	f.Chmod(0000)
	defer func() {
		os.Remove("testdata/test_wrong_permissions")
	}()

	err := Copy("testdata/test_wrong_permissions", "testdata/test_wrong_permissions_copy", 0, 1000, true)
	require.Truef(t, errors.Is(err, os.ErrPermission), "actual error %q", err)
}

func TestInvalidParams(t *testing.T) {
	os.Mkdir("testdata/restricted", 0000)
	defer func() {
		os.Remove("testdata/copied.txt")
		os.Remove("testdata/restricted")
	}()
	testCases := []struct {
		testName      string
		fromPath      string
		toPath        string
		offset        int64
		limit         int64
		expectedError error
	}{
		{
			testName:      "Negative offset",
			fromPath:      "testdata/input.txt",
			toPath:        "testdata/copied.txt",
			offset:        -100,
			limit:         100,
			expectedError: ErrInvalidOffset,
		},
		{
			testName:      "Offset exceeds file size",
			fromPath:      "testdata/input.txt",
			toPath:        "testdata/copied.txt",
			offset:        10000000000000000,
			limit:         100,
			expectedError: ErrOffsetExceedsFileSize,
		},
		{
			testName:      "Negative limit",
			fromPath:      "testdata/input.txt",
			toPath:        "testdata/copied.txt",
			offset:        0,
			limit:         -1000,
			expectedError: ErrInvalidLimit,
		},
		{
			testName:      "File not exists",
			fromPath:      "testdata/not_existed.txt",
			toPath:        "testdata/copied.txt",
			offset:        0,
			limit:         1000,
			expectedError: os.ErrNotExist,
		},
		{
			testName:      "Restricted output",
			fromPath:      "testdata/input.txt",
			toPath:        "testdata/restricted/copy.txt",
			offset:        0,
			limit:         1000,
			expectedError: os.ErrPermission,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.testName, func(t *testing.T) {
			err := Copy(tc.fromPath, tc.toPath, tc.offset, tc.limit, true)
			require.Truef(t, errors.Is(err, tc.expectedError), "actual error %q", err)
		})
	}
}

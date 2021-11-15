package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type flags struct {
	from, to, expected string
	limit, offset      int64
}

func TestCopy(t *testing.T) {
	Flags := map[string]flags{
		"offset exceeds file size": {
			from: "testdata/input.txt", to: "testdata/out.txt",
			offset: 7000,
		},
		"empty": {
			from: "testdata/input_empty.txt", to: "testdata/out.txt",
			expected: "testdata/out_empty.txt", limit: 0, offset: 0,
		},
		"no file": {
			from: "testdata/no_file.txt", to: "testdata/out.txt",
			limit: 0, offset: 0,
		},
		"unknown length": {
			from: "/dev/urandom", to: "testdata/out.txt",
			limit: 0, offset: 0,
		},
		"no path from": {from: "", to: "testdata/out.txt"},
		"no path to":   {from: "testdata/input.txt", to: ""},
		"unsupported file": {
			from: "testdata/input_unsupported.txt", to: "testdata/out.txt",
			limit: 0, offset: 0,
		},
	}

	t.Run("offset exceeds file size", func(t *testing.T) {
		fl := Flags["offset exceeds file size"]
		err := Copy(fl.from, fl.to, fl.offset, fl.limit)

		require.Truef(t, errors.Is(err, ErrOffsetExceedsFileSize), "actual err - %v", err)
	})

	t.Run("empty file", func(t *testing.T) {
		fl := Flags["empty"]
		err := Copy(fl.from, fl.to, fl.offset, fl.limit)
		defer func() {
			err = os.Remove(fl.to)
			if err != nil {
				log.Fatal(err)
			}
		}()

		require.Truef(t, errors.Is(err, ErrEmptyOrUnknownLength), "actual err - %v", err)
		require.True(t, compareFiles(fl.expected, fl.to), "files are not equal")
	})

	t.Run("no file", func(t *testing.T) {
		fl := Flags["no file"]
		err := Copy(fl.from, fl.to, fl.offset, fl.limit)

		require.Truef(t, errors.Is(err, ErrNoSuchFile), "actual err - %v", err)
	})

	t.Run("unknown length", func(t *testing.T) {
		fl := Flags["unknown length"]
		err := Copy(fl.from, fl.to, fl.offset, fl.limit)
		defer func() {
			err = os.Remove(fl.to)
			if err != nil {
				log.Fatal(err)
			}
		}()
		require.True(t, compareFiles("out_empty.txt", fl.to), "files are not equal")
		require.Truef(t, errors.Is(err, ErrEmptyOrUnknownLength), "actual err - %v", err)
	})

	t.Run("no path from", func(t *testing.T) {
		fl := Flags["no path from"]
		err := Copy(fl.from, fl.to, fl.offset, fl.limit)

		require.Truef(t, errors.Is(err, ErrNoPathFrom), "actual err - %v", err)
	})

	t.Run("no path to", func(t *testing.T) {
		fl := Flags["no path to"]
		err := Copy(fl.from, fl.to, fl.offset, fl.limit)

		require.Truef(t, errors.Is(err, ErrNoPathTo), "actual err - %v", err)
	})

	t.Run("unsupported file", func(t *testing.T) {
		fl := Flags["unsupported file"]
		err := Copy(fl.from, fl.to, fl.offset, fl.limit)

		require.Truef(t, errors.Is(err, ErrUnsupportedFile), "actual err - %v", err)
	})
}

func compareFiles(path1, path2 string) bool {
	file1, err := os.Open(path1)
	defer func() {
		err = file1.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()
	if err != nil {
		log.Println(err)
	}

	file2, err := os.Open(path2)
	defer func() {
		err = file2.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()
	if err != nil {
		log.Println(err)
	}

	scan1 := bufio.NewScanner(file1)
	scan2 := bufio.NewScanner(file2)

	for scan1.Scan() {
		scan2.Scan()
		if !bytes.Equal(scan1.Bytes(), scan2.Bytes()) {
			return false
		}
	}

	return true
}

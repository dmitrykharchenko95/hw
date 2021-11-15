package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")

	ErrEmptyOrUnknownLength = errors.New("empty file or unknown length of file")
	ErrNoSuchFile           = errors.New("no such file")
	ErrNoPathFrom           = errors.New("enter path from")
	ErrNoPathTo             = errors.New("enter path to")
)

// Copy копирует limit байт или все данные до io.EOF из fromPath в toPath с отступом offset байт в fromPath.
// Если fromPath пустой файл или файл непределенного размера, Copy возвращает ошибку "empty file or unknown length
// of file" и создает пустой файл toPath.
func Copy(fromPath, toPath string, offset, limit int64) error {
	switch {
	case fromPath == "":
		return ErrNoPathFrom
	case toPath == "":
		return ErrNoPathTo
	}

	fromFile, err := os.OpenFile(fromPath, os.O_RDWR, os.ModeAppend)
	if err != nil {
		switch {
		case err.Error() == "open "+fromPath+": no such file or directory":
			return ErrNoSuchFile
		default:
			return ErrUnsupportedFile
		}
	}
	defer func() {
		err = fromFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	infoFrom, err := fromFile.Stat()

	switch {
	case err != nil:
		return err
	case infoFrom.Size() < offset:
		return ErrOffsetExceedsFileSize
	case infoFrom.Size() == 0:
		toFile, err := os.Create(toPath)
		if err != nil {
			return err
		}
		defer func() {
			err = toFile.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()
		fmt.Printf("\r[%s] 100%%\n", strings.Repeat("=", 100))
		return ErrEmptyOrUnknownLength
	case limit == 0 || limit > infoFrom.Size()-offset:
		limit = infoFrom.Size() - offset
	}

	toFile, err := os.Create(toPath)
	defer func() {
		err = toFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	if err != nil {
		return err
	}

	_, err = fromFile.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	ticker := time.Tick(1 * time.Millisecond)
	var wg sync.WaitGroup
	wg.Add(1)

	go processChecker(toFile, limit, &wg, ticker)
	_, err = io.CopyN(toFile, fromFile, limit)
	if err != nil {
		return err
	}

	wg.Wait()
	return nil
}

func processChecker(toFile *os.File, limit int64, wg *sync.WaitGroup, ticker <-chan time.Time) {
	var (
		infoToFile os.FileInfo
		process    int
	)

	for process != 100 {
		<-ticker
		infoToFile, _ = toFile.Stat()
		process = int(float64(infoToFile.Size())/float64(limit)) * 100
		fmt.Printf("\r[%s%s] %v%%", strings.Repeat("=", process), strings.Repeat("-", 100-process), process)
	}
	fmt.Println()
	wg.Done()
}

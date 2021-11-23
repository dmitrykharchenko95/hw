package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

var (
	ErrWrongFileName  = errors.New("the file name must not contain '='")
	ErrEmptyDirectory = errors.New("directory in empty")
)

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	Env := make(Environment)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, ErrEmptyDirectory
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if strings.Contains(file.Name(), "=") {
			return nil, ErrWrongFileName
		}

		_, ok := os.LookupEnv(file.Name())

		fileData, err := os.Open(dir + "/" + file.Name())
		if err != nil {
			return Env, err
		}

		buf := bufio.NewReader(fileData)
		envValByte, err := buf.ReadBytes(10)
		if err != nil && !errors.Is(err, io.EOF) {
			return Env, err
		}

		envValByte = bytes.ReplaceAll(envValByte, []byte{0}, []byte{10})
		envValString := strings.TrimRightFunc(string(envValByte), func(r rune) bool {
			return r == '\t' || r == ' ' || r == '\n'
		})

		Env[file.Name()] = EnvValue{Value: envValString, NeedRemove: ok}
	}

	return Env, nil
}

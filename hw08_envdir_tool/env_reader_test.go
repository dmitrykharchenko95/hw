package main

import (
	"errors"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	err := os.Mkdir("testdata/additions", 0o777)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer func() {
		err := os.RemoveAll("testdata/additions")
		if err != nil {
			log.Fatal(err)
			return
		}
	}()

	t.Run("name with =", func(t *testing.T) {
		_, err := os.Create("testdata/additions/with=")
		if err != nil {
			log.Fatal(err)
			return
		}
		defer func() {
			err := os.Remove("testdata/additions/with=")
			if err != nil {
				log.Fatal(err)
				return
			}
		}()
		env, err := ReadDir("testdata/additions")

		require.Nil(t, env, "actual value Environment - %v", env)
		require.True(t, errors.Is(err, ErrWrongFileName), "actual err - %v", err)
	})

	t.Run("empty dir", func(t *testing.T) {
		env, err := ReadDir("testdata/additions")

		require.Nil(t, env, "actual value Environment - %v", env)
		require.True(t, errors.Is(err, ErrEmptyDirectory), "actual err - %v", err)
	})

	t.Run("no such dir", func(t *testing.T) {
		env, err := ReadDir("testdata/ghost")

		require.Nil(t, env, "actual value Environment - %v", env)
		require.Error(t, err, "error should not be nil")
	})
}

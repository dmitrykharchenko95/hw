package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	testEnv := Environment{
		"BAR": {Value: "bar", NeedRemove: false},
	}

	t.Run("exit code test", func(t *testing.T) {
		data := "#!/usr/bin/env bash\nexit 77"
		err := ioutil.WriteFile("exit.sh", []byte(data), 0o777)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			err := os.Remove("exit.sh")
			if err != nil {
				log.Fatal(err)
			}
		}()
		exitCode := RunCmd([]string{"./exit.sh"}, testEnv)
		require.True(t, exitCode == 77, "actual ExitCode - %v, expected - 77", exitCode)
	})

	t.Run("nil Environment", func(t *testing.T) {
		exitCode := RunCmd([]string{"ls"}, nil)

		require.True(t, exitCode == 0, "actual - %v", exitCode)
	})
}

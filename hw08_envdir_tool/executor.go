package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	var (
		command   *exec.Cmd
		osCommand = cmd[0]
	)

	if len(cmd) > 1 {
		command = exec.Command(osCommand, cmd[1:]...)
	} else {
		command = exec.Command(osCommand)
	}

	command.Stdout, command.Stdin, command.Stderr = os.Stdout, os.Stdin, os.Stderr

	var appendedEnv []string

	for key, v := range env {
		if v.NeedRemove {
			err := os.Unsetenv(key)
			if err != nil {
				log.Fatalf("os.Unsetenv(%v): %v", key, err)
			}
		}
		if len(v.Value) != 0 {
			appendedEnv = append(appendedEnv, key+"="+v.Value)
		}
	}
	command.Env = append(os.Environ(), appendedEnv...)

	if err := command.Run(); err != nil {
		var e *exec.ExitError
		if errors.As(err, &e) {
			returnCode = command.ProcessState.ExitCode()
		} else {
			log.Fatal(err)
		}
	}
	return
}

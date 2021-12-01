package main

import (
	"os"
)

func main() {
	args := os.Args
	env, _ := ReadDir(args[1])
	RunCmd(args[2:], env)
}

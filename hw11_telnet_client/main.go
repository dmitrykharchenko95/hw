package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

var timeout time.Duration

func init() {
	flag.DurationVar(&timeout, "timeout", time.Second*10, "Timeout connect")
}

func main() {
	flag.Parse()
	args := os.Args
	address := net.JoinHostPort(args[2], args[3])

	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		log.Fatal("Connect err:", err)
	}

	defer func() {
		if err := client.Close(); err != nil {
			log.Fatal("Close err:", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(2)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		if err := client.Send(); err != nil {
			return
		}
		if _, err := fmt.Fprintln(os.Stderr, "...EOF"); err != nil {
			log.Fatal(err)
		}
	}(&wg)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		if err := client.Receive(); err != nil {
			log.Fatal("Receive err:", err)
		}
		if _, err := fmt.Fprintln(os.Stderr, "...Connection was closed by peer"); err != nil {
			log.Fatal(err)
		}
	}(&wg)
	wg.Wait()
}

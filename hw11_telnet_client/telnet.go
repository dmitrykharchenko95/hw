package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	client := &Client{Address: address, Timeout: timeout, In: in, Out: out}
	client.ctx, client.cancel = signal.NotifyContext(context.Background(), syscall.SIGINT)
	return client
}

type Client struct {
	Address string
	Conn    net.Conn
	Timeout time.Duration
	In      io.ReadCloser
	Out     io.Writer

	ctx    context.Context
	cancel context.CancelFunc
}

func (c *Client) Connect() error {
	var err error
	c.Conn, err = net.DialTimeout("tcp", c.Address, c.Timeout)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(os.Stderr, "Connected to %v\n", c.Address)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Close() error {
	if err := c.Conn.Close(); err != nil {
		return err
	}
	return nil
}

func (c *Client) Send() error {
	scanner := bufio.NewScanner(c.In)
loop:
	for {
		select {
		case <-c.ctx.Done():
			break loop
		default:
			if !scanner.Scan() {
				break loop
			}
			if _, err := c.Conn.Write([]byte(fmt.Sprintf("%s\n", scanner.Text()))); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Client) Receive() error {
	scanner := bufio.NewScanner(c.Conn)
loop:
	for {
		select {
		case <-c.ctx.Done():
			break loop
		default:
			if !scanner.Scan() {
				break loop
			}
			_, err := c.Out.Write([]byte(fmt.Sprintf("%s\n", scanner.Text())))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

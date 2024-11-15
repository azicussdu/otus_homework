package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type telnetCli struct {
	address string
	timeout time.Duration
	conn    net.Conn
	in      io.ReadCloser
	out     io.Writer
}

func (tc *telnetCli) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), tc.timeout)
	defer cancel()

	conn, err := (&net.Dialer{}).DialContext(ctx, "tcp", tc.address) // should be for ex: "127.0.0.1:3302"
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	tc.conn = conn

	return nil
}

func (tc *telnetCli) Close() error {
	if tc.conn != nil {
		return tc.conn.Close()
	}
	return nil
}

func (tc *telnetCli) Send() error {
	scanner := bufio.NewScanner(tc.in)
	if tc.conn == nil {
		return fmt.Errorf("connection closed by server")
	}
	for scanner.Scan() {
		text := scanner.Text()

		_, err := tc.conn.Write([]byte(text + "\n"))
		if err != nil {
			return err
		}
	}
	return scanner.Err()
}

func (tc *telnetCli) Receive() error {
	scanner := bufio.NewScanner(tc.conn)
	for scanner.Scan() {
		text := scanner.Text()
		_, err := fmt.Fprintf(tc.out, "%s\n", text)
		if err != nil {
			return fmt.Errorf("failed to receive data: %w", err)
		}
	}
	return scanner.Err()
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	if timeout == 0 {
		timeout = 10 * time.Second // default timeout of 10 seconds
	}
	return &telnetCli{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

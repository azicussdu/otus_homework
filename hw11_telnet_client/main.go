package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go-telnet [--timeout=10s] host port")
		os.Exit(1)
	}

	timeout := 10 * time.Second
	args := os.Args[1:]
	var host, port string

	switch {
	case len(args) == 3 && len(args[0]) > 9 && args[0][:9] == "--timeout":
		var err error
		timeout, err = time.ParseDuration(args[0][10:])
		if err != nil {
			fmt.Println("Invalid timeout format")
			os.Exit(1)
		}
		host, port = args[1], args[2]

	case len(args) == 2:
		host, port = args[0], args[1]

	default:
		fmt.Println("Usage: go-telnet [--timeout=10s] host port")
		os.Exit(1)
	}

	address := host + ":" + port
	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	err := client.Connect()
	if err != nil {
		fmt.Println("Error connecting:", err)
		os.Exit(1)
	}
	defer func(client TelnetClient) {
		err := client.Close()
		if err != nil {
			fmt.Println("Error closing connection:", err)
		}
	}(client)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		select {
		case <-ctx.Done():
			_ = client.Close() // Ensure connection is closed on signal
		}
	}()

	sendErrCh := make(chan error, 1)
	receiveErrCh := make(chan error, 1)

	go func() {
		sendErrCh <- client.Send()
	}()
	go func() {
		receiveErrCh <- client.Receive()
	}()

	select {
	case <-ctx.Done():
		fmt.Println("Received shutdown signal")
	case err := <-sendErrCh:
		if err != nil {
			fmt.Println("Send error:", err)
		}
	case err := <-receiveErrCh:
		if err != nil {
			fmt.Println("Receive error:", err)
		}
	}
}

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run client.go <server-ip>")
		fmt.Println("Example: go run client.go 192.168.1.100")
		os.Exit(1)
	}

	serverAddr := os.Args[1] + ":8080"
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("Error connecting:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Printf("Connected to %s\n", serverAddr)
	fmt.Println("Type messages to send (Ctrl+C to exit):")

	// Read responses in a separate goroutine
	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			fmt.Printf("Echo: %s\n", scanner.Text())
		}
	}()

	// Send messages
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		msg := scanner.Text()
		_, err := conn.Write([]byte(msg + "\n"))
		if err != nil {
			fmt.Println("Error sending:", err)
			return
		}
	}
}

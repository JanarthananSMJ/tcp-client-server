package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

// handleConnection manages the communication with a single client.
func handleConnection(conn net.Conn) {
	// Always close the connection when this function finishes.
	defer conn.Close()

	fmt.Printf("New client connected: %s\n", conn.RemoteAddr())

	// Create a new scanner to read from the connection line by line.
	scanner := bufio.NewScanner(conn)
	for {
		// Read a line from the client. This will block until a newline is sent.
		if scanner.Scan() {
			message := scanner.Text()
			fmt.Printf("Received from %s: %s\n", conn.RemoteAddr(), message)

			// If the client sends "exit", close the connection.
			if strings.ToLower(message) == "exit" {
				fmt.Printf("Client %s disconnected.\n", conn.RemoteAddr())
				break
			}

			// Send a response back to the client.
			response := strings.ToUpper(message)
			fmt.Fprintf(conn, "Server says: %s\n", response)
		} else {
			// If scanning fails (e.g., client disconnects abruptly), break the loop.
			fmt.Printf("Client %s disconnected.\n", conn.RemoteAddr())
			break
		}
	}
}

func main() {
	// 1. Listen for incoming TCP connections on port 8080.
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	// Always close the listener when the main function exits.
	defer listener.Close()

	fmt.Println("Server is listening on :8080...")

	// 2. The main server loop: continuously accept new connections.
	for {
		// Accept() blocks until a new client connects.
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue // Go back to waiting for the next connection.
		}

		// 3. Handle the new connection in a new goroutine.
		// This is the magic! It allows the server to handle multiple clients concurrently.
		go handleConnection(conn)
	}
}

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// IMPORTANT: If the server is on a different machine on your local network,
	// replace "localhost:8080" with the server's local IP address, e.g., "192.168.1.10:8080".
	serverAddress := "localhost:8080"

	// 1. Connect to the server.
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	// Always close the connection when the main function exits.
	defer conn.Close()

	fmt.Printf("Connected to server at %s\n", serverAddress)
	fmt.Println("Type messages and press Enter. Type 'exit' to quit.")

	// Create a reader to get input from the user's terminal.
	reader := bufio.NewReader(os.Stdin)

	for {
		// Read input from the user.
		fmt.Print("> ")
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			break
		}

		// 2. Send the message to the server.
		_, err = fmt.Fprintf(conn, message)
		if err != nil {
			fmt.Println("Error sending message:", err)
			break
		}

		// If the user types 'exit', close the connection.
		if strings.TrimSpace(message) == "exit" {
			break
		}

		// 3. Wait for and read the response from the server.
		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Error reading response:", err)
			break
		}

		fmt.Print("Server response: " + response)
	}

	fmt.Println("Connection closed.")
}

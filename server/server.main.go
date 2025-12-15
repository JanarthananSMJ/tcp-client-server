package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Server listening on 0.0.0.0:8080")
	fmt.Println("Find your IP with: ifconfig (Linux/Mac) or ipconfig (Windows)")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("Client connected: %s\n", conn.RemoteAddr().String())

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		msg := scanner.Text()
		fmt.Printf("Received: %s\n", msg)

		// Echo back
		_, err := conn.Write([]byte(msg + "\n"))
		if err != nil {
			fmt.Println("Error writing:", err)
			return
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Connection error:", err)
	}
	fmt.Printf("Client disconnected: %s\n", conn.RemoteAddr().String())
}

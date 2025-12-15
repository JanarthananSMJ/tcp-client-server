package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
)

type Client struct {
	conn net.Conn
	id   int
}

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Server listening on 0.0.0.0:8080")
	fmt.Println("Waiting for clients to connect...")
	fmt.Println("Press Enter when all clients are connected")

	clients := []Client{}
	clientID := 1
	var mu sync.Mutex

	// Accept clients in background
	done := make(chan bool)
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			mu.Lock()
			clients = append(clients, Client{conn: conn, id: clientID})
			fmt.Printf("Client %d connected from %s\n", clientID, conn.RemoteAddr())
			clientID++
			mu.Unlock()
		}
	}()

	// Wait for user to press Enter
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	close(done)

	mu.Lock()
	numClients := len(clients)
	mu.Unlock()

	if numClients == 0 {
		fmt.Println("No clients connected. Exiting.")
		return
	}

	fmt.Printf("\n%d clients connected\n", numClients)

	// Get range from user
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter start range: ")
	startStr, _ := reader.ReadString('\n')
	start, _ := strconv.Atoi(startStr[:len(startStr)-1])

	fmt.Print("Enter end range: ")
	endStr, _ := reader.ReadString('\n')
	end, _ := strconv.Atoi(endStr[:len(endStr)-1])

	if start > end {
		fmt.Println("Invalid range")
		return
	}

	// Divide work
	totalNums := end - start + 1
	numsPerClient := totalNums / numClients
	remainder := totalNums % numClients

	fmt.Printf("\nDividing range %d-%d among %d clients\n", start, end, numClients)

	results := make([]int, numClients)
	var wg sync.WaitGroup

	currentStart := start
	for i, client := range clients {
		clientNums := numsPerClient
		if i < remainder {
			clientNums++
		}
		clientEnd := currentStart + clientNums - 1

		wg.Add(1)
		go func(c Client, s, e, idx int) {
			defer wg.Done()
			results[idx] = processClient(c, s, e)
		}(client, currentStart, clientEnd, i)

		fmt.Printf("Client %d: range %d-%d (%d numbers)\n", client.id, currentStart, clientEnd, clientNums)
		currentStart = clientEnd + 1
	}

	fmt.Println("\nWaiting for client results...\n")
	wg.Wait()

	// Calculate final sum
	totalSum := 0
	fmt.Println("Results received from clients:")
	for i, result := range results {
		fmt.Printf("Client %d returned: %d\n", clients[i].id, result)
		totalSum += result
	}

	fmt.Printf("\n=== FINAL SUM OF DIGITS: %d ===\n", totalSum)

	// Close all connections
	for _, client := range clients {
		client.conn.Close()
	}
}

func processClient(c Client, start, end int) int {
	// Send client ID and range
	msg := fmt.Sprintf("%d %d %d\n", c.id, start, end)
	_, err := c.conn.Write([]byte(msg))
	if err != nil {
		fmt.Printf("Error sending to client %d: %v\n", c.id, err)
		return 0
	}

	scanner := bufio.NewScanner(c.conn)
	if scanner.Scan() {
		result, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Printf("Error parsing result from client %d: %v\n", c.id, err)
			return 0
		}
		fmt.Printf("✓ Received result from Client %d: %d\n", c.id, result)
		return result
	}

	fmt.Printf("Client %d disconnected\n", c.id)
	return 0
}

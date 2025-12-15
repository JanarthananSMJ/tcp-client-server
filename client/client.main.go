package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run client.go <server-ip>")
		fmt.Println("Example: go run client.go 192.168.1.41")
		fmt.Println("Or use 'localhost' for local connection")
		os.Exit(1)
	}

	serverAddr := os.Args[1] + ":8080"
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("Error connecting:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Printf("Connected to server at %s\n", serverAddr)
	fmt.Println("Waiting for work assignment...")

	scanner := bufio.NewScanner(conn)
	if !scanner.Scan() {
		fmt.Println("Connection closed by server")
		return
	}

	// Parse client ID and range from server
	parts := strings.Split(scanner.Text(), " ")
	if len(parts) != 3 {
		fmt.Println("Invalid message from server")
		return
	}

	clientID, _ := strconv.Atoi(parts[0])
	start, _ := strconv.Atoi(parts[1])
	end, _ := strconv.Atoi(parts[2])

	fmt.Printf("\n=== CLIENT %d ===\n", clientID)
	fmt.Printf("Assigned range: %d to %d\n", start, end)

	// Calculate sum of digits
	result := calculateSumOfDigits(start, end)

	fmt.Printf("\nCalculation complete!\n")
	fmt.Printf("Client %d result: %d\n", clientID, result)
	fmt.Println("==================")

	// Send result back
	_, err = conn.Write([]byte(fmt.Sprintf("%d\n", result)))
	if err != nil {
		fmt.Println("Error sending result:", err)
		return
	}

	fmt.Println("\nResult sent to server. Done.")
}

func calculateSumOfDigits(start, end int) int {
	total := 0
	fmt.Println("\nProcessing numbers:")
	for num := start; num <= end; num++ {
		digitSum := sumDigits(num)
		fmt.Printf("  %d → digit sum = %d\n", num, digitSum)
		total += digitSum
	}
	return total
}

func sumDigits(n int) int {
	sum := 0
	if n < 0 {
		n = -n
	}
	for n > 0 {
		sum += n % 10
		n /= 10
	}
	return sum
}

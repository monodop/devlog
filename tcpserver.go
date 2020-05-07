package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func startTcpListener(exitChannel chan bool, messageChannel chan string) {
	address := ":9090"

	listener, err := net.Listen("tcp4", address)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()

	nextId := 1
	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		id := nextId
		nextId++
		go handleConnection(connection, id, messageChannel)
	}
}

func handleConnection(connection net.Conn, id int, messageChannel chan string) {
	fmt.Printf("Opened TCP connection %d to %s\n", id, connection.RemoteAddr().String())
	defer fmt.Printf("Closed TCP connection %d to %s\n", id, connection.RemoteAddr().String())
	for {
		data, err := bufio.NewReader(connection).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		line := strings.TrimSpace(string(data))
		fmt.Printf("%d: %s\n", id, line)
		messageChannel <- line

		connection.Write([]byte("Thanks\n"))
	}
	// connection.Close()
}

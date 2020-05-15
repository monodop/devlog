package main

import (
	"bufio"
	"net"
	"strings"

	"github.com/monodop/devlog/log"
)

func startTcpListener(exitChannel chan bool, messageChannel chan string) {
	address := ":9090"

	listener, err := net.Listen("tcp4", address)
	if err != nil {
		log.Exception(err)
		return
	}
	defer listener.Close()
	log.Info("TCP server now listening on %s", address)

	nextId := 1
	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Exception(err)
			return
		}

		id := nextId
		nextId++
		go handleConnection(connection, id, messageChannel)
	}
}

func handleConnection(connection net.Conn, id int, messageChannel chan string) {
	log.Info("Opened TCP connection %d to %s", id, connection.RemoteAddr().String())
	defer log.Info("Closed TCP connection %d to %s", id, connection.RemoteAddr().String())
	for {
		data, err := bufio.NewReader(connection).ReadString('\n')
		if err != nil {
			log.Exception(err)
			return
		}

		line := strings.TrimSpace(string(data))
		log.Info("%d: %s", id, line)
		messageChannel <- line

		connection.Write([]byte("Thanks\n"))
	}
	// connection.Close()
}

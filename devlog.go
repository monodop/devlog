package main

import (
	"github.com/monodop/devlog/log"
)

func main() {
	log.Info("DevLog starting up...")
	exitChannel := make(chan bool)
	messageChannel := make(chan string)
	go startWsListener(exitChannel, messageChannel)
	go startTcpListener(exitChannel, messageChannel)
	// go startWsTestConnection()
	// go startWsTestConnection()
	<-exitChannel
}

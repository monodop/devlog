package main

import "fmt"

func main() {
	fmt.Println("Hello, world")
	exitChannel := make(chan bool)
	messageChannel := make(chan string)
	go startWsListener(exitChannel, messageChannel)
	go startTcpListener(exitChannel, messageChannel)
	// go startWsTestConnection()
	// go startWsTestConnection()
	<-exitChannel
}

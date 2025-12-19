package main

import "fmt"

type MpesaEnvironment int

const (
	Production MpesaEnvironment = iota
	Sandbox
)

func main() {
	messages := make(chan string)

	var sender chan<- string = messages
	var receiver <-chan string = messages

	// Send in goroutine
	go func() {
		sender <- "Hello"
		sender <- "World"
		sender <- "!!"
	}()

	// Receive in main
	msg1 := <-receiver
	msg2 := <-receiver
	msg3 := <-receiver

	fmt.Println(msg1, msg2, msg3) // Hello World
}

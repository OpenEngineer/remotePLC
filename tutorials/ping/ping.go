package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	// client implementation
	conn, err := net.Dial("udp", "192.168.1.1:7")

	if err != nil {
		log.Fatal(err)
	}

	a := byte(0xaa)
	fmt.Println("sending message...")
	n, err := conn.Write([]byte{a})
	if n <= 0 || err != nil {
		log.Fatal("error sending message")
	}
	fmt.Println("sent message of length ", n)

	buffer := make([]byte, 1024)
	fmt.Println("waiting for response...")
	conn.Read(buffer)
	fmt.Println("response received")
	fmt.Println(string(buffer))
}

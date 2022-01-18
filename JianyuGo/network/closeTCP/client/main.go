package main

import (
	"log"
	"net"
)

func main() {
	// Part 1: open a TCP session to server
	c, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatalf("Error to open TCP connection: %s", err)
	}
	defer c.Close()

	// Part2: write some data to server
	log.Printf("TCP session open")
	b := []byte("Hi, gopher?")
	_, err = c.Write(b)
	if err != nil {
		log.Fatalf("Error writing TCP session: %s", err)
	}

	// Part3: read any responses until get an error
	for {
		d := make([]byte, 100)
		_, err := c.Read(d)
		if err != nil {
			log.Fatalf("Error reading TCP session: %s", err)
		}
		log.Printf("reading data from server: %s\n", string(d))
	}
}

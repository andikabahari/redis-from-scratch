package main

import (
	"io"
	"log"
	"net"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		log.Fatal("Failed to bind to port 6379")
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal("Error accepting connection: ", err.Error())
		}
		defer conn.Close()

		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	buf := make([]byte, 512)
	for {
		var err error

		_, err = conn.Read(buf)
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatal("Error reading data: ", err.Error())
		}

		_, err = conn.Write([]byte("+PONG\r\n"))
		if err != nil {
			log.Fatal("Error writing response: ", err.Error())
		}
	}
}

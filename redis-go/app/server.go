package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
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

var cached map[string]string = make(map[string]string)

func set(key, value string) {
	cached[key] = value
}

func get(key string) string {
	return cached[key]
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

		args := strings.Split(string(buf), "\r\n")

		command := args[2]
		var msg string
		switch command {
		case "ping":
			msg += "+PONG\r\n"
		case "echo":
			text := args[4]
			msg += fmt.Sprintf("+%s\r\n", text)
		case "set":
			key := args[4]
			value := args[6]
			set(key, value)
			msg += "+OK\r\n"
		case "get":
			key := args[4]
			value := get(key)
			msg += fmt.Sprintf("+%s\r\n", value)
		default:
			msg += fmt.Sprintf("-ERR unknown command '%s'\r\n", command)
		}

		_, err = conn.Write([]byte(msg))
		if err != nil {
			log.Fatal("Error writing response: ", err.Error())
		}
	}
}

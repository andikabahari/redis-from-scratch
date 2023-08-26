package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
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

var cached = make(map[string]value)

type value struct {
	value     string
	expiredAt time.Time
}

func set(k, v string, px int64) {
	var exp time.Time
	if px > 0 {
		now := time.Now().UnixMilli()
		exp = time.UnixMilli(px + now)
	}
	value := value{
		value:     v,
		expiredAt: exp,
	}
	cached[k] = value
}

func get(k string) (string, bool) {
	v, ok := cached[k]
	if !ok {
		return "", false
	}
	if !v.expiredAt.IsZero() {
		if time.Now().After(v.expiredAt) {
			delete(cached, k)
			return "", false
		}
	}
	return v.value, true
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
			var px int64
			if len(args) > 8 {
				command = args[8]
				if command == "px" {
					exp := args[10]
					px, err = strconv.ParseInt(exp, 10, 64)
					if err != nil {
						log.Println("Error parsing int: ", err.Error())
					}
				}
			}
			key := args[4]
			value := args[6]
			set(key, value, px)
			msg += "+OK\r\n"
		case "get":
			key := args[4]
			value, ok := get(key)
			if ok {
				msg += fmt.Sprintf("+%s\r\n", value)
			} else {
				msg += "*-1\r\n"
			}
		default:
			msg += fmt.Sprintf("-ERR unknown command '%s'\r\n", command)
		}

		_, err = conn.Write([]byte(msg))
		if err != nil {
			log.Fatal("Error writing response: ", err.Error())
		}
	}
}

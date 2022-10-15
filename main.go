package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"strings"
	"sync"
)

const (
	HOST = "0.0.0.0"
	PORT = ":9001"
	TYPE = "tcp"
)

var DB map[string]string
var mut sync.RWMutex
var once sync.Once

func createServer() error {
	listener, err := net.Listen(TYPE, HOST+PORT)
	if err != nil {
		return err
	}

	defer listener.Close()
	fmt.Println("Listening!")
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go handleIncoming(conn)
	}
}

func handleIncoming(conn net.Conn) {
	reader := bufio.NewReader(conn)
	var buffer bytes.Buffer
	for {
		n, isPrefix, err := reader.ReadLine()
		if err != nil {
			fmt.Printf("%+v\n", err)
			return
		}
		buffer.Write(n)

		if buffer.Len() > 0 && !isPrefix {
			command := strings.Split(buffer.String(), " ")
			switch command[0] {
			case "QUIT":
				conn.Write([]byte("Bye\n"))
				conn.Close()
				return
			case "SET":
				mut.Lock()
				DB[command[1]] = strings.Join(command[2:], " ")
				mut.Unlock()
				conn.Write([]byte("OK\n"))
			case "GET":
				mut.RLock()
				if val, ok := DB[command[1]]; ok {
					conn.Write([]byte("OK " + val + "\n"))
				} else {
					conn.Write([]byte("OK\n"))
				}
				mut.RUnlock()
			default:
				conn.Write([]byte("OK unimplemented\n"))
			}
			buffer.Reset()
		}

	}
}

func init() {
	once.Do(func() {
		DB = make(map[string]string)
	})
}

func main() {
	createServer()
}

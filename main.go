package main

import (
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

	return nil
}

func handleIncoming(conn net.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Printf("%+v\n", err)
			return
		}

		if n > 0 {
			data := strings.TrimSpace(string(buf[:n]))
			command := strings.Split(data, " ")
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

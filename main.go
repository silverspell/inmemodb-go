package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"sync"

	"net/http"
	_ "net/http/pprof"
)

const (
	HOST = "0.0.0.0"
	PORT = ":9001"
	TYPE = "tcp"
)

var DB map[string][]byte
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
	buffer := bytes.Buffer{}
	for {
		n, isPrefix, err := reader.ReadLine()
		if err != nil {
			fmt.Printf("%+v\n", err)
			return
		}

		buffer.Write(n)

		if buffer.Len() > 0 && !isPrefix {
			command := buffer.Bytes()[0:4]
			switch string(command) {
			case "QUIT":
				conn.Write([]byte("Bye\n"))
				conn.Close()
				return
			case "SET ":
				key, val, _ := keyValSeperator(buffer.Bytes()[4:])
				mut.Lock()
				DB[key] = val
				mut.Unlock()
				conn.Write([]byte("OK\n"))
			case "GET ":
				key, _, _ := keyValSeperator(buffer.Bytes()[4:])
				mut.RLock()
				if val, ok := DB[key]; ok {
					response := []byte("OK ")
					response = append(response, val...)
					response = append(response, []byte("\n")...)
					conn.Write(response)
				} else {
					conn.Write([]byte("OK\n"))
				}
				mut.RUnlock()
			case "PING":
				conn.Write([]byte("OK PONG\n"))
			default:
				conn.Write([]byte("OK unimplemented\n"))
			}
			buffer.Reset()
		}

	}
}

func keyValSeperator(buffer []byte) (string, []byte, error) {
	for i, val := range buffer {
		if val == 32 { // 32 = " "
			return string(buffer[0:i]), buffer[i+1:], nil
		}
	}
	return string(buffer), nil, nil
}

func init() {
	once.Do(func() {
		DB = make(map[string][]byte)
	})
}

func main() {
	go func() {
		http.ListenAndServe("0.0.0.0:8080", nil)
	}()
	createServer()
}

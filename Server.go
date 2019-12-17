package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

const (
	protocol string = "tcp"
	// host     string = "127.0.0.1"
	port     string = ":81"
)

func main() {
	// boot up the tcp server
	listener, err := net.Listen(protocol, port)
	if err != nil {
		log.Fatal("tcp server failiure: ", err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("TCP server accept error:", err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	bufferBytes, err := bufio.NewReader(conn).ReadBytes('\n')

	if err != nil {
		log.Println("Connection broken by client..")
		conn.Close()
		return
	}

	message := string(bufferBytes)
	clientAddress := conn.RemoteAddr().String()
	response := fmt.Sprintf(message + "from " + clientAddress + "\n")
	log.Println(response)
	conn.Write([]byte("you sent: " + response))

	handleConnection(conn)
}

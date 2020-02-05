package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
)

const (
	protocol string = "tcp"
	// uncomment if you need the server for localhost otherwise don't specify if you want remote connections
	// host     string = "127.0.0.1"
	port string = ":81"
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
	// read buffer from client after connection is established
	bufferBytes, err := bufio.NewReader(conn).ReadBytes('\n')

	if err != nil {
		log.Println("Connection broken by client, ", err)
		conn.Close()
		// escape recursion
		return
	}

	// convert the buffer bytes to string data type
	message := string(bufferBytes)

	// get client's IP address
	clientAddress := conn.RemoteAddr().String()
	response := fmt.Sprintf(message + "from " + clientAddress + "\n")

	// check key word to allow data in that connection to processed
	// example data you might receive from a sensor maybe KEYWORD,(1,1,1,1,2016-02-05,):
	if strings.Contains(message, "KEYWORD") {
		important := message[7 : len(message)-1]
		result := strings.SplitAfter(important, ",")
		// first := strings.Replace(result[0], "(", "", -1)
		user := strings.Replace(result[3], ",", "", -1)
		location := strings.Replace(result[3], ",", "", -1)
		finger := strings.Replace(result[2], ",", "", -1)
		device := strings.Replace(result[3], ",", "", -1)
		action := strings.Replace(result[3], ",", "", -1)
		requestBody, err := json.Marshal(map[string]string{
			"guard":        user,
			"guard_finger": finger,
			"location":     location,
			"device":       device,
			"action":       action,
		})

		if err != nil {
			log.Fatalln(err)
		}
		// make a http post request to your endpoint
		enpointURL := "http://0a58d63a.ngrok.io/api/v1/data-stream/"

		resp, err := http.Post(enpointURL, "application/json", bytes.NewReader(requestBody))
		//  error handling
		if err != nil {
			log.Fatalln(err)
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			log.Fatalln(err)
		}
		// print the response body to the server's terminal
		log.Println(string(body))
	}
	log.Println(response)

	// let the client know what happened
	conn.Write([]byte("you sent: " + response))

	// recursive func to handle io.EOF for random disconnects
	handleConnection(conn)
}

// http://0a58d63a.ngrok.io/api/v1/data-stream/

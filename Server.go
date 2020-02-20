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
	"strconv"
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
	
	 type Data struct {
	    action int64  `json:"ref"`
	    device int64  `json:"ref"`
	    guard_finger int64  `json:"ref"`
	    tag int64  `json:"ref"`
	    flag_one int64  `json:"ref"`
	    flag_two int64  `json:"ref"`
	    arrived_at string
  	}

	// get client's IP address
	clientAddress := conn.RemoteAddr().String()
	response := fmt.Sprintf(message + "from " + clientAddress + "\n")

	// check key word to allow data in that connection to processed
	// example data you might receive from a sensor maybe KEYWORD,255,7,,10,0,0,2016-02-05 02:59:20 ,:
	if strings.Contains(message, "PITJET") {
		important := message[7 : len(message)-1]
		result := strings.SplitAfter(important, ",")
		device := strings.Replace(result[0], ",", "", -1)
		finger := strings.Replace(result[1], ",", "", -1)
		tag := strings.Replace(result[2], ",", "", -1)
		action := strings.Replace(result[3], ",", "", -1)
		flag1 := strings.Replace(result[4], ",", "", -1)
		flag2 := strings.Replace(result[5], ",", "", -1)
		datetime := strings.Replace(result[6], ",", "", -1)
		
		device, err := strconv.ParseInt(device,10, 64)
		finger, err := strconv.ParseInt(finger,10, 64)
		tag, err := strconv.ParseInt(tag,10, 64)
		action := strconv.ParseInt(action,10, 64)
		flag1, err := strconv.ParseInt(flag1,10, 64)
		flag2, err := strconv.ParseInt(flag2,10, 64)
		
		data := Data{
			device: device,
			guard_finger: finger,
			flag_one: flag1,
			flag_two: flag2,
			tag: tag,
			action: action,
			arrived_at: datetime,
		}
		
		var requestBody []byte
		requestBody, err := json.Marshal(data)

		if err != nil {
			log.Fatalln(err)
		}
		// make a http post request to your endpoint
		enpointURL := "http://httpbin.org/post"

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

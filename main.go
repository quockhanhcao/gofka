package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	// "strconv"
)

func main() {
	if os.Args[1] == "server" {
		// startServer()
		broker := Broker{}
		err := broker.startBrokerServer()
		if err != nil {
			fmt.Printf("error starting broker: %v\n", err.Error())
		}
	} else if os.Args[1] == "producer" {
		// fmt.Println("starting producer")
		// port, err := strconv.ParseInt(os.Args[2], 10, 32)
		// if err != nil {
		// 	panic(err)
		// }
		// producer := Producer{}

	} else {
		clientConnectAndEcho(10000)
	}
}

func clientConnectAndEcho(port int) {
	conn, err := net.Dial("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Printf("error %v", err)
		return
	}
	// read input stdin, write to stream
	streamReadWrite := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return
	}

	msg := Message{
		ECHO: &line,
	}

	fmt.Printf("Sending msg to server %s", *msg.ECHO)
	err = writeMessageToStream(streamReadWrite, msg)
	if err != nil {
		panic(err)
	}

	// read back from the stream
	resp, err := readMessageFromStream(streamReadWrite)
	fmt.Printf("Received msg from server %s\n", *resp.R_ECHO)
	conn.Close()
}

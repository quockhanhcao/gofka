package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
)

func main() {
	switch os.Args[1] {
	case "server":
		broker := Broker{}
		err := broker.startBrokerServer()
		if err != nil {
			fmt.Printf("error starting broker: %v\n", err.Error())
		}
	case "producer":
		port, err := strconv.ParseInt(os.Args[2], 10, 32)
		if err != nil {
			panic(err)
		}
		producer := Producer{}
		producer.startProducerServer(int16(port))

	default:
		clientConnectAndEcho(10000)
	}
}

func clientConnectAndEcho(port int) {
	conn, err := net.Dial("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return
	}
	// read input stdin, write to stream
	stream_rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return
	}

	msg := Message{
		ECHO: &line,
	}

	fmt.Printf("Sending msg to server %s", *msg.ECHO)
	err = writeMessageToStream(stream_rw, msg)
	if err != nil {
		panic(err)
	}

	// read back from the stream
	resp, err := readMessageFromStream(stream_rw)
	fmt.Printf("Received msg from server %s\n", *resp.R_ECHO)
	conn.Close()
}

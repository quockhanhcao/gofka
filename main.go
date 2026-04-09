package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	if os.Args[1] == "server" {
		// startServer()
		broker := Broker{}
		err := broker.startBrokerServer()
		if err != nil {
			return
		}
	} else {
		clientConnect()
	}
}

func startServer() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Printf("error creating server %v", err)
		return
	}
	conn, err := ln.Accept()
	if err != nil {
		fmt.Printf("error accepting connection %v", err)
		return
	}
	streamReadWrite := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	length, _ := streamReadWrite.ReadByte()
	data, _ := streamReadWrite.Peek(int(length))
	streamReadWrite.Discard(int(length))

	fmt.Printf("received from client %s", string(data))

	streamReadWrite.WriteByte(byte(len(data)))
	streamReadWrite.WriteString(string(data))
	streamReadWrite.Flush()

	conn.Close()
}

func clientConnect() {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		fmt.Printf("error %v", err)
		return
	}
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	fmt.Printf("sending to server %s", line)
	// send to server
	// a simple protocol: a message has to attach a byte for message length before the actual message
	// streamWriter.WriteByte(byte(len(line)))
	// streamWriter.WriteString(line)
	// streamWriter.Flush()

	streamWriter := bufio.NewWriter(conn)
	if strings.HasPrefix(line, "ECHO") {
		streamWriter.WriteByte(byte(len(line[5:]) + 1))
		streamWriter.WriteByte(1)
		streamWriter.WriteString(line[5:])
		// fmt.Printf("actual data %s", line[4:])
		streamWriter.Flush()
	}

	// var str []byte
	// str = append(str, byte(len(line)))
	// str = append(str, []byte(line)...)

	// conn.Write(str)

	streamReader := bufio.NewReader(conn)
	length, _ := streamReader.ReadByte()
	data, _ := streamReader.Peek(int(length))

	fmt.Printf("received from server: %s\n", string(data))
	conn.Close()
}

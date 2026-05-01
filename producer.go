package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

type Producer struct {
}

func (p *Producer) registerBroker(port int16) error {
	conn, err := net.Dial("tcp", fmt.Sprintf(":%d", BROKER_PORT))
	if err != nil {
		return err
	}
	// read input stdin, write to stream
	stream_rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	portStr := fmt.Sprintf("%d", port)
	msg := Message{
		P_REG: &portStr,
	}
	err = writeMessageToStream(stream_rw, msg)
	if err != nil {
		panic(err)
	}

	// read back from the stream
	resp, err := readMessageFromStream(stream_rw)
	fmt.Printf("Received response from broker %v\n", *resp.R_P_REG)
	return nil
}

func (p *Producer) startProducerServer(port int16) error {
	var err error

	// connect to broker to send register
	err = p.registerBroker(port)
	if err != nil {
		return err
	}

	ln, err := net.Listen(PROTOCOL, fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	conn, err := ln.Accept()
	if err != nil {
		return err
	}
	stream_rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}

		fmt.Printf("/////////////// Sending msg to broker: %s\n", line)
		// write ECHO to stream
		err = writeMessageToStream(stream_rw, Message{
			ECHO: &line,
		})
		if err != nil {
			break
		}

		// read response
		resp, err := readMessageFromStream(stream_rw)
		if err != nil {
			break
		}
		fmt.Printf("Received msg from broker: %s\n", *resp.R_ECHO)
	}

	// close connection
	err = conn.Close()
	if err != nil {
		return err
	}
	return nil
}

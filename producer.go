package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

type Producer struct {
	port    uint16
	topicID uint16
}

func (p *Producer) registerBroker() error {
	conn, err := net.Dial("tcp", fmt.Sprintf(":%d", BROKER_PORT))
	if err != nil {
		return err
	}
	// read input stdin, write to stream
	stream_rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	msg := Message{
		P_REG: &ProducerRegisterMessage{
			port:    p.port,
			topicID: p.topicID,
		},
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

func (p *Producer) startProducerServer() error {
	var err error

	// connect to broker to send register
	err = p.registerBroker()
	if err != nil {
		return err
	}

	ln, err := net.Listen(PROTOCOL, fmt.Sprintf(":%d", p.port))
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

		// write ECHO to stream
		err = writeMessageToStream(stream_rw, Message{
			PCM: []byte(line),
		})
		if err != nil {
			break
		}

		// read response
		resp, err := readMessageFromStream(stream_rw)
		if err != nil {
			break
		}
		fmt.Printf("Received msg from broker: %v\n", *resp.R_PCM)
	}

	// close connection
	err = conn.Close()
	if err != nil {
		return err
	}
	return nil
}

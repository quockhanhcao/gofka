package main

import (
	"bufio"
	"fmt"
	"net"
)

const BROKER_PORT = 10000
const PROTOCOL = "tcp"

type Broker struct {
}

func (b *Broker) startBrokerServer() error {
	ln, err := net.Listen(PROTOCOL, fmt.Sprintf(":%d", BROKER_PORT))
	if err != nil {
		// fmt.Printf("error creating server %v", err)
		return err
	}
	for {
		var err error
		conn, err := ln.Accept()
		if err != nil {
			// fmt.Printf("error accepting connection %v", err)
			return err
		}
		fmt.Println("Accepting connection")
		streamReadWrite := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
		parseMsg, err := readMessageFromStream(streamReadWrite)
		if err != nil {
			return err
		}

		if parseMsg != nil {
			resp, err := b.processBrokerMessage(parseMsg)
			if err != nil {
				return err
			}
			err = writeMessageToStream(streamReadWrite, *resp)
			if err != nil {
				return err
			}
		}

		// close connection
		err = conn.Close()
		if err != nil {
			return err
		}
	}
}

// process message, call inner process function
// response correct message
func (b *Broker) processBrokerMessage(msg *Message) (*Message, error) {
	if msg.ECHO != nil {
		resp, err := b.processEchoMessage(msg.ECHO)
		if err != nil {
			return nil, err
		}
		return &Message{R_ECHO: &resp}, nil
	}

	return nil, nil
}

func (b *Broker) processEchoMessage(msg *string) (string, error) {
	return fmt.Sprintf("I have received: %s", *msg), nil
}

package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"time"
)

const BROKER_PORT = 10000
const PROTOCOL = "tcp"

type Broker struct {
}

func (b *Broker) startBrokerServer() error {
	ln, err := net.Listen(PROTOCOL, fmt.Sprintf(":%d", BROKER_PORT))
	if err != nil {
		return err
	}
	for {
		var err error
		conn, err := ln.Accept()
		if err != nil {
			// return err
			break
		}
		stream_rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
		parsedMsg, err := readMessageFromStream(stream_rw)
		if err != nil {
			fmt.Println("Error here after reading msg")
			// return err
			break
		}

		if err == nil && parsedMsg != nil {
			resp, err := b.processBrokerMessage(parsedMsg)
			if err != nil {
				fmt.Println("Error here after process msg")
				// return err
				break
			}
			err = writeMessageToStream(stream_rw, *resp)
			if err != nil {
				fmt.Println("Error here after write msg")
				// return err
				break
			}
		}

		// close connection
		// err = conn.Close()
		// if err != nil {
		// 	fmt.Println("Error here after closing connection")
		// 	// return err
		// 	break
		// }
		// fmt.Println("Already close connection")
	}
	return nil
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

	if msg.P_REG != nil {
		resp, err := b.processProducerRegMsg(msg.P_REG)
		if err != nil {
			return nil, err
		}
		return &Message{R_P_REG: resp}, nil
	}

	return nil, nil
}

func (b *Broker) processEchoMessage(msg *string) (string, error) {
	return fmt.Sprintf("I have received: %s", *msg), nil
}

func (b *Broker) processProducerRegMsg(msg *string) (*byte, error) {
	port, err := strconv.ParseInt(*msg, 10, 32)
	if err != nil {
		return nil, err
	}
	go func() {
		var conn net.Conn
		var err error
		for i := 0; i < 10; i++ {
			conn, err = net.Dial(PROTOCOL, fmt.Sprintf(":%d", port))
			if err == nil {
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
		if err != nil {
			fmt.Println("failed to connect to producer:", err)
			return
		}
		stream_rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
		for {
			msg, err := readMessageFromStream(stream_rw)
			if msg == nil && err != nil {
				// panic(err)
				fmt.Println("error read msg from stream")
				break
			}
			// process msg
			resp, err := b.processBrokerMessage(msg)
			err = writeMessageToStream(stream_rw, *resp)
			if err != nil {
				// panic(err)
				fmt.Println("error write msg to stream")
				break
			}
		}
	}()

	var resp byte = 0
	return &resp, nil
}

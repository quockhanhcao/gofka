package main

import (
	"bufio"
	"fmt"
	"net"
)

const (
	ECHO = 1
)

func readFromStream(streamReadWrite *bufio.ReadWriter) ([]byte, error) {
	length, err := streamReadWrite.ReadByte()
	if err != nil {
		return nil, err
	}
	data, err := streamReadWrite.Peek(int(length))
	if err != nil {
		return nil, err
	}
	_, err = streamReadWrite.Discard(int(length))
	if err != nil {
		return nil, err
	}

	return data, nil
}

func writeToStream(streamReadWrite *bufio.ReadWriter, data string) error {
	err := streamReadWrite.WriteByte(byte(len(data)))
	if err != nil {
		return err
	}
	_, err = streamReadWrite.WriteString(string(data))
	if err != nil {
		return err
	}
	err = streamReadWrite.Flush()
	if err != nil {
		return err
	}

	return nil
}

type Broker struct {
}

func (b *Broker) startBrokerServer() error {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		// fmt.Printf("error creating server %v", err)
		return err
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// fmt.Printf("error accepting connection %v", err)
			return err
		}
		fmt.Println("Accepting connection")
		streamReadWrite := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
		data, err := readFromStream(streamReadWrite)
		if err != nil {
			return err
		}

		fmt.Printf("received message from client %s", string(data))

		// parse message here
		msg := b.parseMessage(data)
		if msg != nil {
			resp, err := b.processMessage(msg)
			if err != nil {
				return err
			}
			// write back response
			err = writeToStream(streamReadWrite, resp)
			if err != nil {
				return err
			}
		}

		// write it back
		err = conn.Close()
		if err != nil {
			return err
		}
	}
}

type Message struct {
	ECHO *string
}

func (b *Broker) parseMessage(data []byte) *Message {
	switch data[0] {
	case ECHO:
		fmt.Println("echo case")
		fmt.Printf("string data %s", string(data[1:]))
		st := string(data[1:])
		return &Message{
			ECHO: &st,
		}
	default:
		return nil
	}
}

func (b *Broker) processMessage(msg *Message) (string, error) {
	var resp string
	var err error
	if msg.ECHO != nil {
		resp = fmt.Sprintf("I have received: %s", *msg.ECHO)
		return resp, nil
	}
	return resp, err
}
